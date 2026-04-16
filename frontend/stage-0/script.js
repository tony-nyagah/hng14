const DUE = new Date("2026-04-19T23:59:00");

function formatDueDate(date) {
  return (
    "Due " +
    date.toLocaleDateString("en-US", {
      month: "short",
      day: "numeric",
      year: "numeric",
    })
  );
}

function formatTimeRemaining(due) {
  const diffMs = due - Date.now();
  const el = document.getElementById("time-remaining-display");

  if (diffMs < 0) {
    const abs = Math.abs(diffMs);
    const days = Math.floor(abs / 86400000);
    const hrs = Math.floor(abs / 3600000);
    const mins = Math.floor(abs / 60000);
    el.style.color = "#ff3b3b";
    if (days >= 1) return `Overdue by ${days} day${days !== 1 ? "s" : ""}`;
    if (hrs >= 1) return `Overdue by ${hrs} hour${hrs !== 1 ? "s" : ""}`;
    return `Overdue by ${mins} minute${mins !== 1 ? "s" : ""}`;
  }

  const days = Math.floor(diffMs / 86400000);
  const hrs = Math.floor(diffMs / 3600000);
  const mins = Math.floor(diffMs / 60000);

  if (hrs < 6) {
    el.style.color = "#ff3b3b";
  } else if (days < 2) {
    el.style.color = "#ff8c00";
  } else {
    el.style.color = "#007700";
  }

  if (days >= 1) return `Due in ${days} day${days !== 1 ? "s" : ""}`;
  if (hrs >= 1) return `Due in ${hrs} hour${hrs !== 1 ? "s" : ""}`;
  return `Due in ${mins} minute${mins !== 1 ? "s" : ""}`;
}

function updateTimes() {
  document.getElementById("due-date-display").textContent = formatDueDate(DUE);
  document.getElementById("time-remaining-display").textContent =
    formatTimeRemaining(DUE);
}

updateTimes();
setInterval(updateTimes, 30000);
