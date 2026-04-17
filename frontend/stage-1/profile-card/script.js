(function () {
  'use strict';

  const timeEl = document.getElementById('epoch-time');

  function updateTime() {
    const ms = Date.now();
    timeEl.textContent = ms;
    // Keep the datetime attribute in sync for accessibility
    timeEl.setAttribute('datetime', new Date(ms).toISOString());
  }

  // Set immediately on load, then refresh every 1000ms
  updateTime();
  setInterval(updateTime, 1000);
})();
