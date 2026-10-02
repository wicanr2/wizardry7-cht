// Disposable UI prototype only. No original game, recordings or network inputs.
const puppeteer = require('puppeteer');
const fs = require('node:fs');
const crypto = require('node:crypto');
const assert = require('node:assert/strict');

(async () => {
  const browser = await puppeteer.launch({args: ['--no-sandbox', '--disable-dev-shm-usage']});
  try {
    const page = await browser.newPage();
    await page.setViewport({width: 1200, height: 850});
    const errors = [], externalRequests = [];
    page.on('pageerror', error => errors.push(error.message));
    page.on('request', request => { if (!request.url().startsWith('file:')) externalRequests.push(request.url()); });
    await page.goto('file:///src/recording-stop-prototype.html');
    const phase = id => page.$eval(`#${id}`, card => card.dataset.phase);
    const value = (id, field) => page.$eval(`#${id} .${field}`, node => node.textContent);
    const waitFor = (id, state) => page.waitForFunction((a, b) => document.getElementById(a).dataset.phase === b, {}, id, state);
    const cases = [];
    await page.screenshot({path: '/out/comparison.png'});
    assert.equal(await phase('limit'), 'idle');
    assert.equal(await phase('manual'), 'idle');
    cases.push('initial-idle');
    await page.click('#limit .start');
    for (let i = 0; i < 2; i++) await page.click('#limit .advance');
    assert.equal(await phase('limit'), 'recording');
    assert.equal(await value('limit', 'step'), '400000');
    await page.click('#limit .advance');
    await waitFor('limit', 'finished');
    assert.equal(await value('limit', 'step'), '600000');
    assert.equal(await value('limit', 'count'), '3');
    await page.click('#limit .advance');
    assert.equal(await value('limit', 'step'), '600000');
    cases.push('limit-stops-at-target-and-rejects-further-input');
    await page.click('#limit .reset');
    await page.$eval('#limit .target', input => { input.value = '300000'; });
    await page.click('#limit .start');
    assert.equal(await phase('limit'), 'idle');
    cases.push('invalid-target-rejected');
    await page.click('#limit .reset');
    await page.$eval('#limit .target', input => { input.value = '600000'; });
    await page.click('#manual .start');
    await page.click('#manual .advance');
    await page.click('#manual .stop');
    assert.equal(await phase('manual'), 'waiting');
    assert.equal(await page.$eval('#manual .advance', node => node.disabled), true);
    await page.$eval('#manual .advance', node => node.dispatchEvent(new Event('click')));
    assert.equal(await value('manual', 'step'), '200000');
    await page.screenshot({path: '/out/manual-waiting.png'});
    await waitFor('manual', 'finished');
    assert.equal(await value('manual', 'count'), '1');
    cases.push('manual-stops-at-boundary-and-rejects-new-input');
    await page.click('#manual .reset');
    await page.click('#manual .start');
    await page.click('#manual .stop');
    await page.click('#manual .reset');
    await page.click('#manual .start');
    await page.click('#manual .advance');
    await new Promise(resolve => setTimeout(resolve, 600));
    assert.equal(await phase('manual'), 'recording');
    assert.equal(await value('manual', 'step'), '200000');
    cases.push('reset-invalidates-prior-stop-timer');
    await page.setViewport({width: 390, height: 844});
    const fits = await page.evaluate(() => document.documentElement.scrollWidth <= innerWidth);
    assert.equal(fits, true);
    cases.push('narrow-screen-without-horizontal-overflow');
    assert.deepEqual(errors, []);
    assert.deepEqual(externalRequests, []);
    const hash = path => crypto.createHash('sha256').update(fs.readFileSync(path)).digest('hex');
    const report = {scope: 'disposable UI simulation; no game or recording integration',
      node: process.version, browser: await browser.version(), cases,
      errors, external_requests: externalRequests,
      source_sha256: hash('/src/recording-stop-prototype.html'),
      test_sha256: hash('/src/recording-stop-prototype-test.cjs'),
      comparison_sha256: hash('/out/comparison.png'), waiting_sha256: hash('/out/manual-waiting.png')};
    fs.writeFileSync('/out/report.json', JSON.stringify(report, null, 2) + '\n');
    console.log(JSON.stringify(report));
  } finally {
    await browser.close();
  }
})().catch(error => { console.error(error); process.exitCode = 1; });
