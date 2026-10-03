// Disposable UI simulation. No game or recording inputs.
const puppeteer=require('puppeteer');
const fs=require('node:fs');
const crypto=require('node:crypto');
const assert=require('node:assert/strict');
(async()=>{
  const browser=await puppeteer.launch({args:['--no-sandbox','--disable-dev-shm-usage']});
  try {
    const page=await browser.newPage();await page.setViewport({width:1200,height:900});
    const errors=[],requests=[],cases=[];
    page.on('pageerror',e=>errors.push(e.message));
    page.on('request',r=>{if(!r.url().startsWith('file:'))requests.push(r.url());});
    await page.goto('file:///src/recording-stop-continue-prototype.html');
    const value=(id,cls)=>page.$eval(`#${id} .${cls}`,n=>n.textContent);
    const phase=id=>page.$eval(`#${id}`,n=>n.dataset.phase);
    const wait=(id,state)=>page.waitForFunction((a,b)=>document.getElementById(a).dataset.phase===b,{},id,state);
    const click=(id,cls)=>page.click(`#${id} .${cls}`);
    const force=(id,cls)=>page.$eval(`#${id} .${cls}`,n=>n.dispatchEvent(new Event('click')));
    const expect=async(id,expected)=>{for(const [k,v] of Object.entries(expected))assert.equal(await value(id,k),v);};
    assert.equal(await phase('exit'),'idle');assert.equal(await phase('continue'),'idle');
    await page.screenshot({path:'/out/comparison.png',fullPage:true});cases.push('initial-idle');
    await click('exit','start');await click('exit','advance');await click('exit','stop');
    assert.equal(await phase('exit'),'waiting');await force('exit','advance');
    await wait('exit','finished');await force('exit','advance');await force('exit','start');
    await expect('exit',{step:'200000',total:'1',count:'1',end:'200000'});
    cases.push('exit-finishes-and-rejects-later-operations');
    await click('continue','start');await click('continue','advance');await click('continue','advance');await click('continue','stop');
    assert.equal(await phase('continue'),'waiting');await force('continue','advance');await force('continue','stop');
    await expect('continue',{step:'400000',total:'2',count:'2',end:'尚未停止'});
    await page.screenshot({path:'/out/continue-waiting.png',fullPage:true});
    await wait('continue','finalizing');await force('continue','advance');
    await expect('continue',{step:'400000',total:'2',count:'2',end:'400000'});
    await wait('continue','playing');cases.push('continue-blocks-input-through-boundary-and-finalization');
    for(let i=0;i<3;i++)await click('continue','advance');
    await force('continue','stop');await force('continue','start');
    assert.equal(await phase('continue'),'playing');
    await expect('continue',{step:'1000000',total:'5',count:'2',end:'400000'});
    await page.screenshot({path:'/out/continue-playing.png',fullPage:true});
    cases.push('continued-play-advances-game-only-and-keeps-recording-cutoff');
    await click('continue','reset');await click('continue','start');await click('continue','stop');
    await click('continue','reset');await click('continue','start');await click('continue','advance');
    await new Promise(r=>setTimeout(r,650));assert.equal(await phase('continue'),'recording');
    await expect('continue',{step:'200000',total:'1',count:'1',end:'尚未停止'});
    cases.push('reset-cancels-old-boundary-timer');
    await click('continue','stop');await wait('continue','finalizing');
    await click('continue','reset');await click('continue','start');await click('continue','advance');
    await new Promise(r=>setTimeout(r,650));assert.equal(await phase('continue'),'recording');
    await expect('continue',{step:'200000',total:'1',count:'1',end:'尚未停止'});
    cases.push('reset-cancels-old-finalization-timer');
    await page.setViewport({width:390,height:844});
    assert.equal(await page.evaluate(()=>document.documentElement.scrollWidth<=innerWidth),true);
    cases.push('narrow-screen-without-horizontal-overflow');
    assert.deepEqual(errors,[]);assert.deepEqual(requests,[]);
    const hash=p=>crypto.createHash('sha256').update(fs.readFileSync(p)).digest('hex');
    const report={scope:'disposable UI only; no engine, recording writer or audio integration',node:process.version,browser:await browser.version(),cases,errors,external_requests:requests,
      source_sha256:hash('/src/recording-stop-continue-prototype.html'),test_sha256:hash('/src/recording-stop-continue-prototype-test.cjs'),
      images:Object.fromEntries(['comparison.png','continue-waiting.png','continue-playing.png'].map(n=>[n,hash('/out/'+n)]))};
    fs.writeFileSync('/out/report.json',JSON.stringify(report,null,2)+'\n');console.log(JSON.stringify(report));
  }finally{await browser.close();}
})().catch(e=>{console.error(e);process.exitCode=1;});
