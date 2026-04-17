<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, nextTick } from 'vue'

type Priority = 'Low' | 'Medium' | 'High'
type Status = 'Pending' | 'In Progress' | 'Done'

// Task state
const title = ref('Design the Landing Page')
const description = ref(
  'Create wireframes and high-fidelity mockups for the new product landing page. Coordinate with the marketing team for final copy and imagery. Review competitor pages and gather inspiration for the layout and visual direction.',
)
const priority = ref<Priority>('High')
const dueDate = ref('2026-04-19')
const status = ref<Status>('In Progress')

// Edit mode state
const isEditing = ref(false)
const editTitle = ref('')
const editDescription = ref('')
const editPriority = ref<Priority>('High')
const editDueDate = ref('')
const editButtonRef = ref<HTMLButtonElement | null>(null)

// Expand / collapse
const COLLAPSE_THRESHOLD = 120
const isExpanded = ref(false)
const shouldCollapse = computed(() => description.value.length > COLLAPSE_THRESHOLD)

// Tick every 30 s
const now = ref(Date.now())
let timer: ReturnType<typeof setInterval>

onMounted(() => {
  timer = setInterval(() => {
    now.value = Date.now()
  }, 30000)
})
onUnmounted(() => clearInterval(timer))

// Checkbox ↔ status two-way binding
const isDone = computed({
  get: () => status.value === 'Done',
  set: (val: boolean) => {
    status.value = val ? 'Done' : 'Pending'
  },
})

// Parsed due-date object
const dueDateObj = computed(() => {
  if (!dueDate.value) return null
  return new Date(dueDate.value + 'T23:59:00')
})

const formattedDueDate = computed(() => {
  if (!dueDateObj.value) return ''
  return (
    'Due ' +
    dueDateObj.value.toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
      year: 'numeric',
    })
  )
})

const isOverdue = computed(() => {
  if (!dueDateObj.value || status.value === 'Done') return false
  return dueDateObj.value.getTime() < now.value
})

const timeRemaining = computed(() => {
  if (status.value === 'Done') return 'Completed'
  if (!dueDateObj.value) return ''
  const diff = dueDateObj.value.getTime() - now.value
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
  if (status.value === 'Done') return '#16a34a'
  const diff = dueDateObj.value ? dueDateObj.value.getTime() - now.value : 0
  const h = diff / 3600000
  const d = diff / 86400000
  if (diff < 0 || h < 6) return '#ef4444'
  if (d < 2) return '#f97316'
  return '#16a34a'
})

// Priority accent colour
const priorityColor = computed(() => {
  switch (priority.value) {
    case 'High':
      return '#ef4444'
    case 'Medium':
      return '#f97316'
    case 'Low':
      return '#22c55e'
  }
})

// Status badge background
const statusBadgeBg = computed(() => {
  switch (status.value) {
    case 'Done':
      return '#b2ff59'
    case 'In Progress':
      return '#ffe500'
    case 'Pending':
      return '#e0e0e0'
  }
})

// Edit actions
function startEdit() {
  editTitle.value = title.value
  editDescription.value = description.value
  editPriority.value = priority.value
  editDueDate.value = dueDate.value
  isEditing.value = true
}

async function saveEdit() {
  title.value = editTitle.value
  description.value = editDescription.value
  priority.value = editPriority.value
  dueDate.value = editDueDate.value
  isEditing.value = false
  await nextTick()
  editButtonRef.value?.focus()
}

async function cancelEdit() {
  isEditing.value = false
  await nextTick()
  editButtonRef.value?.focus()
}
</script>

<template>
  <main class="min-h-screen bg-[#f0f0f0] flex items-center justify-center p-8">
    <article data-testid="test-todo-card" class="neo-card bg-white max-w-[480px] w-full"
      :class="{ 'opacity-70': status === 'Done' }">
      <!-- Priority indicator strip (top accent bar) -->
      <div data-testid="test-todo-priority-indicator" class="h-2 w-full" :style="{ background: priorityColor }"
        :aria-label="`Priority indicator: ${priority}`" />

      <div class="p-6">
        <!-- ─── EDIT FORM ─── -->
        <form v-if="isEditing" data-testid="test-todo-edit-form" class="flex flex-col gap-4" @submit.prevent="saveEdit">
          <div class="flex flex-col gap-1">
            <label for="edit-title" class="meta-label">Title</label>
            <input id="edit-title" data-testid="test-todo-edit-title-input" v-model="editTitle" type="text" required
              class="neo-input" />
          </div>

          <div class="flex flex-col gap-1">
            <label for="edit-description" class="meta-label">Description</label>
            <textarea id="edit-description" data-testid="test-todo-edit-description-input" v-model="editDescription"
              rows="4" class="neo-input resize-none" />
          </div>

          <div class="grid grid-cols-2 gap-3">
            <div class="flex flex-col gap-1">
              <label for="edit-priority" class="meta-label">Priority</label>
              <select id="edit-priority" data-testid="test-todo-edit-priority-select" v-model="editPriority"
                class="neo-input">
                <option value="Low">Low</option>
                <option value="Medium">Medium</option>
                <option value="High">High</option>
              </select>
            </div>
            <div class="flex flex-col gap-1">
              <label for="edit-due-date" class="meta-label">Due date</label>
              <input id="edit-due-date" data-testid="test-todo-edit-due-date-input" v-model="editDueDate" type="date"
                class="neo-input" />
            </div>
          </div>

          <div class="flex gap-3 pt-2">
            <button type="submit" data-testid="test-todo-save-button" class="neo-btn bg-[#b2ff59] text-black flex-1"
              @click="saveEdit">
              Save
            </button>
            <button type="button" data-testid="test-todo-cancel-button" class="neo-btn bg-white text-black flex-1"
              @click="cancelEdit">
              Cancel
            </button>
          </div>
        </form>

        <!-- ─── CARD VIEW ─── -->
        <template v-else>
          <!-- Header: title + priority badge -->
          <div class="flex items-start justify-between gap-4 mb-3">
            <h2 data-testid="test-todo-title" class="text-xl font-black leading-tight flex-1"
              :class="{ 'line-through text-gray-400': status === 'Done' }">
              {{ title }}
            </h2>
            <span data-testid="test-todo-priority" class="neo-badge text-white" :style="{ background: priorityColor }"
              :aria-label="`Priority: ${priority}`">
              {{ priority }}
            </span>
          </div>

          <!-- Description (collapsible) -->
          <div id="collapsible-desc" data-testid="test-todo-collapsible-section"
            class="mb-1 overflow-hidden transition-[max-height] duration-300"
            :class="shouldCollapse && !isExpanded ? 'max-h-[4.5rem]' : 'max-h-[999px]'">
            <p data-testid="test-todo-description" class="text-sm leading-relaxed">
              {{ description }}
            </p>
          </div>

          <button v-if="shouldCollapse" data-testid="test-todo-expand-toggle" type="button"
            class="text-xs font-black uppercase tracking-wide mb-5 underline cursor-pointer bg-transparent border-none p-0"
            :aria-expanded="isExpanded" aria-controls="collapsible-desc" @click="isExpanded = !isExpanded">
            {{ isExpanded ? 'Show less ▲' : 'Show more ▼' }}
          </button>
          <div v-else class="mb-5" />

          <hr class="border-t-2 border-black mb-4" />

          <!-- Due date + time remaining -->
          <div class="grid grid-cols-2 gap-3 mb-4 text-sm">
            <div class="flex flex-col gap-0.5">
              <span class="meta-label">Due date</span>
              <time data-testid="test-todo-due-date" :datetime="dueDate" class="font-bold"
                :class="{ 'text-[#ef4444]': isOverdue }">
                {{ formattedDueDate }}
              </time>
            </div>
            <div class="flex flex-col gap-0.5">
              <span class="meta-label">Time remaining</span>
              <time data-testid="test-todo-time-remaining" :datetime="dueDate" class="font-black"
                :style="{ color: timeColor }">
                {{ timeRemaining }}
              </time>
            </div>
          </div>

          <!-- Overdue indicator (only shown when overdue) -->
          <div v-if="isOverdue" data-testid="test-todo-overdue-indicator" role="alert"
            class="mb-4 inline-flex items-center gap-1.5 px-3 py-1 border-2 border-[#ef4444] text-[#ef4444] text-xs font-black uppercase tracking-wide">
            <span aria-hidden="true">⚠</span>
            Overdue
          </div>

          <!-- Status display + control + checkbox -->
          <div class="flex items-center gap-3 flex-wrap mb-5">
            <span data-testid="test-todo-status" class="neo-badge text-black" :style="{ background: statusBadgeBg }"
              :aria-label="`Status: ${status}`">
              {{ status }}
            </span>

            <select data-testid="test-todo-status-control" v-model="status" class="neo-select text-xs"
              aria-label="Change status">
              <option value="Pending">Pending</option>
              <option value="In Progress">In Progress</option>
              <option value="Done">Done</option>
            </select>

            <label class="flex items-center gap-2 cursor-pointer ml-auto">
              <input type="checkbox" id="complete-toggle" data-testid="test-todo-complete-toggle" v-model="isDone"
                aria-label="Mark task as complete" class="neo-checkbox" />
              <span class="text-xs font-black uppercase tracking-wide select-none">Done</span>
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
            <button ref="editButtonRef" data-testid="test-todo-edit-button" aria-label="Edit task" type="button"
              class="neo-btn bg-[#ffe500] text-black" @click="startEdit">
              Edit
            </button>
            <button data-testid="test-todo-delete-button" aria-label="Delete task" type="button"
              class="neo-btn bg-[#ff3b3b] text-white">
              Delete
            </button>
          </div>
        </template>
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
  transition:
    box-shadow 0.08s,
    transform 0.08s;
}

.neo-btn:hover {
  box-shadow: 5px 5px 0 #000;
  transform: translate(-1px, -1px);
}

.neo-btn:active {
  box-shadow: 1px 1px 0 #000;
  transform: translate(2px, 2px);
}

.neo-input {
  font-family: 'Arial Black', Arial, sans-serif;
  font-size: 0.85rem;
  font-weight: 700;
  padding: 8px 10px;
  border: 2px solid #000;
  border-radius: 0;
  background: #fff;
  width: 100%;
  outline: none;
  box-shadow: 3px 3px 0 #000;
}

.neo-input:focus {
  box-shadow: 5px 5px 0 #000;
}

.neo-select {
  font-family: 'Arial Black', Arial, sans-serif;
  font-weight: 700;
  padding: 5px 8px;
  border: 2px solid #000;
  border-radius: 0;
  background: #fff;
  cursor: pointer;
  box-shadow: 2px 2px 0 #000;
}
</style>
