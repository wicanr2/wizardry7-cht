"""容器內FreeType的唯讀樣張工具；ABI依映像內freetype.h，僅Linux64。"""
import ctypes as c
import ctypes.util


class Generic(c.Structure):
    _fields_ = [('data',c.c_void_p),('finalizer',c.c_void_p)]


class Box(c.Structure):
    _fields_ = [(s,c.c_long) for s in ('xmin','ymin','xmax','ymax')]


class Vector(c.Structure):
    _fields_ = [('x',c.c_long),('y',c.c_long)]


class Metrics(c.Structure):
    _fields_ = [(s,c.c_long) for s in ('width','height','hbx','hby','ha','vbx','vby','va')]


class Bitmap(c.Structure):
    _fields_ = [('rows',c.c_uint),('width',c.c_uint),('pitch',c.c_int),('buffer',c.c_void_p),
                ('num_grays',c.c_ushort),('pixel_mode',c.c_ubyte),('palette_mode',c.c_ubyte),('palette',c.c_void_p)]


class Slot(c.Structure):
    _fields_ = [('library',c.c_void_p),('face',c.c_void_p),('next',c.c_void_p),('glyph_index',c.c_uint),
                ('generic',Generic),('metrics',Metrics),('linear_h',c.c_long),('linear_v',c.c_long),
                ('advance',Vector),('format',c.c_uint),('bitmap',Bitmap),('left',c.c_int),('top',c.c_int)]


class Face(c.Structure):
    _fields_ = [(s,c.c_long) for s in ('num_faces','face_index','face_flags','style_flags','num_glyphs')] + [
        ('family_name',c.c_char_p),('style_name',c.c_char_p),('num_fixed_sizes',c.c_int),('available_sizes',c.c_void_p),
        ('num_charmaps',c.c_int),('charmaps',c.c_void_p),('generic',Generic),('bbox',Box),
        ('units',c.c_ushort)] + [(s,c.c_short) for s in ('ascender','descender','height','max_aw','max_ah','underline_position','underline_thickness')] + [
        ('glyph',c.POINTER(Slot)),('size',c.c_void_p),('charmap',c.c_void_p)]


class Font:
    def __init__(self,path,index=2):
        assert c.sizeof(c.c_void_p)==8 and c.sizeof(c.c_long)==8, '僅支援已核Linux64 ABI'
        self.ft = c.CDLL(ctypes.util.find_library('freetype'))
        self.lib = c.c_void_p()
        self.face = c.POINTER(Face)()
        self.ft.FT_Init_FreeType.argtypes = [c.POINTER(c.c_void_p)]
        self.ft.FT_New_Face.argtypes = [c.c_void_p,c.c_char_p,c.c_long,c.POINTER(c.POINTER(Face))]
        self.ft.FT_Set_Pixel_Sizes.argtypes = [c.POINTER(Face),c.c_uint,c.c_uint]
        self.ft.FT_Get_Char_Index.argtypes = [c.POINTER(Face),c.c_ulong]
        self.ft.FT_Get_Char_Index.restype = c.c_uint
        self.ft.FT_Load_Char.argtypes = [c.POINTER(Face),c.c_ulong,c.c_int32]
        self.ft.FT_Done_Face.argtypes = [c.POINTER(Face)]
        self.ft.FT_Done_FreeType.argtypes = [c.c_void_p]
        self.ft.FT_Library_Version.argtypes = [c.c_void_p,c.POINTER(c.c_int),c.POINTER(c.c_int),c.POINTER(c.c_int)]
        assert self.ft.FT_Init_FreeType(c.byref(self.lib))==0
        assert self.ft.FT_New_Face(self.lib,str(path).encode(),index,c.byref(self.face))==0
        self.family = self.face.contents.family_name.decode()
        self.style = self.face.contents.style_name.decode()
        a,b,d = c.c_int(),c.c_int(),c.c_int()
        self.ft.FT_Library_Version(self.lib,c.byref(a),c.byref(b),c.byref(d))
        self.version = f'{a.value}.{b.value}.{d.value}'

    def glyph(self,char,size):
        assert self.ft.FT_Get_Char_Index(self.face,ord(char))!=0, f'缺字U+{ord(char):04X}'
        assert self.ft.FT_Set_Pixel_Sizes(self.face,0,size)==0
        assert self.ft.FT_Load_Char(self.face,ord(char),4)==0
        s = self.face.contents.glyph.contents
        b = s.bitmap
        assert b.pixel_mode==2 and b.num_grays==256 and 0<b.width<=128 and 0<b.rows<=128
        stride = abs(b.pitch)
        assert stride>=b.width
        raw = c.string_at(b.buffer,stride*b.rows)
        alpha = []
        for y in range(b.rows):
            row = y if b.pitch>=0 else b.rows-1-y
            alpha.extend(raw[row*stride:row*stride+b.width])
        return dict(w=b.width,h=b.rows,left=s.left,top=s.top,advance=(s.advance.x+32)//64,alpha=alpha)

    def close(self):
        if self.face:
            self.ft.FT_Done_Face(self.face);self.face = c.POINTER(Face)()
        if self.lib:
            self.ft.FT_Done_FreeType(self.lib);self.lib = c.c_void_p()
