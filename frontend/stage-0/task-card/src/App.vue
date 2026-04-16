<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'

const DUE = new Date('2026-04-19T23:59:00')

const isDone = ref(false)
const now = ref(Date.now())

let timer: ReturnType<typeof setInterval>

onMounted(() => {
  timer = setInterval(() => {
    now.value = Date.now()
  }, 30000)
})

onUnmounted(() => clearInterval(timer))

const formattedDueDate = computed(() => {
  return (
    'Due ' +
    DUE.toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })
  )
})

const timeRemaining = computed(() => {
  const diff = DUE.getTime() - now.value
  if (diff < 0) {
    const abs = Math.abs(diff)
    const d = Math.floor(abs / 86400000)
    const h = Math.floor(abs / 3600000)
    const m = Math.floor(abs / 60000)
    if (d >= 1) return `Overdue by ${d} day${d !== 1 ? 's' : ''}`
    if (h >= 1) return `Overdue by ${h} hour${h !== 1 ? 's' : ''}`
    return `Overdue by ${m} minute${m !== 1 ? 's' : ''}`
  }
  const d = Math.floor(diff / 86400000)
  const h = Math.floor(diff / 3600000)
  const m = Math.floor(diff / 60000)
  if (d >= 1) return `Due in ${d} day${d !== 1 ? 's' : ''}`
  if (h >= 1) return `Due in ${h} hour${h !== 1 ? 's' : ''}`
  return `Due in ${m} minute${m !== 1 ? 's' : ''}`
})

const timeColor = computed(() => {
  const diff = DUE.getTime() - now.value
  const h = diff / 3600000
  const d = diff / 86400000
  if (diff < 0 || h < 6) return '#ef4444'
  if (d < 2) return '#f97316'
  return '#16a34a'
})
</script>

<template>
  <main class="min-h-screen bg-[#f0f0f0] flex items-center justify-center p-8">
    <article data-testid="test-todo-card" class="neo-card bg-white max-w-[480px] w-full p-6">

      <!-- Header: title + priority -->
      <div class="flex items-start justify-between gap-4 mb-3">
        <h2 data-testid="test-todo-title" class="text-xl font-black leading-tight flex-1">
          Design the Landing Page
        </h2>
        <span data-testid="test-todo-priority" class="neo-badge bg-[#ff3b3b] text-white"
          aria-label="Priority: High">High</span>
      </div>

      <!-- Description -->
      <p data-testid="test-todo-description" class="text-sm leading-relaxed mb-5">
        Create wireframes and high-fidelity mockups for the new product landing
        page. Coordinate with the marketing team for final copy and imagery.
      </p>

      <hr class="border-t-2 border-black mb-4" />

      <!-- Due date + time remaining -->
      <div class="grid grid-cols-2 gap-3 mb-5 text-sm">
        <div class="flex flex-col gap-0.5">
          <span class="meta-label">Due date</span>
          <time data-testid="test-todo-due-date" datetime="2026-04-19T23:59:00" class="font-bold">{{ formattedDueDate
            }}</time>
        </div>
        <div class="flex flex-col gap-0.5">
          <span class="meta-label">Time remaining</span>
          <time data-testid="test-todo-time-remaining" datetime="2026-04-19T23:59:00" class="font-black"
            :style="{ color: timeColor }">{{ timeRemaining }}</time>
        </div>
      </div>

      <!-- Status + checkbox -->
      <div class="flex items-center gap-4 flex-wrap mb-5">
        <span data-testid="test-todo-status" class="neo-badge bg-[#ffe500] text-black"
          aria-label="Status: In Progress">In Progress</span>

        <label class="flex items-center gap-2 cursor-pointer">
          <input type="checkbox" id="complete-toggle" data-testid="test-todo-complete-toggle" v-model="isDone"
            aria-label="Mark task as complete" class="neo-checkbox" />
          <span class="text-xs font-black uppercase tracking-wide select-none">Mark complete</span>
        </label>
      </div>

      <!-- Tags -->
      <div class="mb-5">
        <span class="meta-label block mb-1.5">Categories</span>
        <ul data-testid="test-todo-tags" role="list" class="flex flex-wrap gap-1.5 list-none p-0 m-0">
          <li data-testid="test-todo-tag-design" class="neo-tag">Design</li>
          <li data-testid="test-todo-tag-ui" class="neo-tag">UI</li>
          <li data-testid="test-todo-tag-frontend" class="neo-tag">Frontend</li>
        </ul>
      </div>

      <!-- Actions -->
      <div class="flex gap-3 pt-4 border-t-2 border-black">
        <button data-testid="test-todo-edit-button" aria-label="Edit task"
          class="neo-btn bg-[#ffe500] text-black">Edit</button>
        <button data-testid="test-todo-delete-button" aria-label="Delete task"
          class="neo-btn bg-[#ff3b3b] text-white">Delete</button>
      </div>

    </article>
  </main>
</template>

<style scoped>
.neo-card {
  border: 3px solid #000;
  box-shadow: 7px 7px 0 #000;
  border-radius: 0;
}

.neo-badge {
  display: inline-block;
  padding: 4px 10px;
  border: 2px solid #000;
  border-radius: 0;
  font-size: 0.65rem;
  font-weight: 900;
  text-transform: uppercase;
  letter-spacing: 0.06em;
  white-space: nowrap;
  flex-shrink: 0;
  align-self: flex-start;
}

.meta-label {
  font-size: 0.6rem;
  font-weight: 900;
  text-transform: uppercase;
  letter-spacing: 0.07em;
  color: #555;
}

.neo-tag {
  font-size: 0.7rem;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  padding: 3px 9px;
  border: 2px solid #000;
  background: #fff;
  box-shadow: 2px 2px 0 #000;
  border-radius: 0;
}

.neo-checkbox {
  width: 20px;
  height: 20px;
  border: 3px solid #000;
  appearance: none;
  -webkit-appearance: none;
  background: #fff;
  cursor: pointer;
  flex-shrink: 0;
  position: relative;
  border-radius: 0;
}

.neo-checkbox:checked {
  background: #b2ff59;
}

.neo-checkbox:checked::after {
  content: '✓';
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  font-size: 13px;
  font-weight: 900;
  color: #000;
  line-height: 1;
}

.neo-btn {
  font-family: 'Arial Black', Arial, sans-serif;
  font-size: 0.78rem;
  font-weight: 900;
  text-transform: uppercase;
  letter-spacing: 0.06em;
  padding: 8px 18px;
  border: 2px solid #000;
  cursor: pointer;
  box-shadow: 3px 3px 0 #000;
  border-radius: 0;
  transition: box-shadow 0.08s, transform 0.08s;
}

.neo-btn:hover {
  box-shadow: 5px 5px 0 #000;
  transform: translate(-1px, -1px);
}

.neo-btn:active {
  box-shadow: 1px 1px 0 #000;
  transform: translate(2px, 2px);
}
</style>
