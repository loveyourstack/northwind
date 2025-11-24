<template>
  <v-sheet height="64">
    <v-toolbar flat>
      <v-btn class="me-4" color="grey-darken-2" variant="outlined" @click="setToday">Today</v-btn>

      <v-btn color="grey-darken-2" size="small" variant="text" icon @click="prev">
        <v-icon size="small">mdi-chevron-left</v-icon>
      </v-btn>

      <v-btn color="grey-darken-2" size="small" variant="text" icon @click="next">
        <v-icon size="small">mdi-chevron-right</v-icon>
      </v-btn>

      <v-toolbar-title v-if="calendar">
        {{ calendar.title }}
      </v-toolbar-title>

      <v-menu location="bottom end">
        <template v-slot:activator="{ props }">
          <v-btn color="grey-darken-2" variant="outlined" v-bind="props">
            <span>{{ getTypeLabel(type) }}</span>
            <v-icon end>mdi-menu-down</v-icon>
          </v-btn>
        </template>
        <v-list>
          <v-list-item @click="type = 'day'">
            <v-list-item-title>Day</v-list-item-title>
          </v-list-item>
          <v-list-item @click="type = 'week'">
            <v-list-item-title>Week</v-list-item-title>
          </v-list-item>
          <v-list-item @click="type = 'month'">
            <v-list-item-title>Month</v-list-item-title>
          </v-list-item>
        </v-list>
      </v-menu>

    </v-toolbar>
  </v-sheet>

  <v-sheet height="700" width="800">

    <v-calendar
      ref="calendar"
      v-model="focus"
      :event-color="getEventColor"
      event-overlap-mode="stack"
      :event-overlap-threshold="30"
      :events="events"
      :type="type"
      color="primary"
      :weekdays="[1, 2, 3, 4, 5]"
      @click:date="viewDay"
      @click:event="showEvent"
      @click:more="viewDay"
    ></v-calendar>

    <v-menu v-if="selectedEvent && selectedElement" v-model="selectedOpen" :activator="selectedElement" :close-on-content-click="false" location="end">
      <v-card color="grey-lighten-4" min-width="350px" flat>
        <v-toolbar :color="getEventColor(selectedEvent)" dark>
          <v-btn icon>
            <v-icon>mdi-pencil</v-icon>
          </v-btn>
          <v-toolbar-title v-html="selectedEvent.type"></v-toolbar-title>
          <v-btn icon>
            <v-icon>mdi-heart</v-icon>
          </v-btn>
          <v-btn icon>
            <v-icon>mdi-dots-vertical</v-icon>
          </v-btn>
        </v-toolbar>
        <v-card-text>
          <span v-html="selectedEvent.name"></span>
        </v-card-text>
        <v-card-actions>
          <v-btn color="secondary" variant="text" @click="selectedOpen = false">
            Cancel
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-menu>

  </v-sheet>
</template>

<script lang="ts" setup>

// adapted from examples: https://vuetifyjs.com/en/components/calendars/#calendars

import { VCalendar } from 'vuetify/labs/VCalendar'
import { ref, onMounted } from 'vue'

interface evt {
  type: string
  name: string | undefined
  start: Date
  end: Date
  timed: boolean
}

const calendar = ref()
const type = ref<any>('month')
const focus = ref('')
const events = ref<evt[]>([])

const selectedEvent = ref<evt>()
const selectedElement = ref(null)
const selectedOpen = ref(false)

function getEventColor(event: any) {
  switch (event.type) {
    case 'Holiday': return 'blue'
    case 'Public holiday': return 'orange'
    case 'Sick': return 'green'
    case 'Training': return 'indigo'

    default: return 'grey'
  }
}

function getTypeLabel(type: string) {
  switch (type) {
    case 'day': return 'Day'
    case 'month': return 'Month'
    case 'week': return 'Week'
  
    default: console.log('unknown type: ' + type)
      break;
  }
}

function next() {
  calendar.value.next()
}
function prev() {
  calendar.value.prev()
}
function setToday() {
  focus.value = ''
}

function showEvent(nativeEvent: any, { event }: { event: any}) {
  const open = () => {
    selectedEvent.value = event
    selectedElement.value = nativeEvent.target
    requestAnimationFrame(() => requestAnimationFrame(() => selectedOpen.value = true))
  }
  if (selectedOpen.value) {
    selectedOpen.value = false
    requestAnimationFrame(() => requestAnimationFrame(() => open()))
  } else {
    open()
  }
  nativeEvent.stopPropagation()
}

function viewDay (nativeEvent: any, { date }: { date: string }) {
  focus.value = date
  type.value = 'day'
}

onMounted(() => {
  events.value.push({
    type: 'Public holiday',
    name: 'Steven, Laura, Anne, Robert',
    start: new Date('2025-11-06'),
    end: new Date('2025-11-06'),
    timed: false,
  })

  events.value.push({
    type: 'Sick',
    name: 'Steven',
    start: new Date('2025-11-19'),
    end: new Date('2025-11-21'),
    timed: false,
  })

  events.value.push({
    type: 'Sick',
    name: 'Andrew',
    start: new Date('2025-11-26T14:00:00'),
    end: new Date('2025-11-26T18:00:00'),
    timed: true,
  })

  events.value.push({
    type: 'Sick',
    name: 'Robert',
    start: new Date('2025-11-26'),
    end: new Date('2025-11-26'),
    timed: false,
  })

  events.value.push({
    type: 'Holiday',
    name: 'Laura',
    start: new Date('2025-11-20'),
    end: new Date('2025-12-05'),
    timed: false,
  })

  events.value.push({
    type: 'Holiday',
    name: 'Nancy',
    start: new Date('2025-11-26'),
    end: new Date('2025-11-26'),
    timed: false,
  })

  events.value.push({
    type: 'Training',
    name: 'Anne',
    start: new Date('2025-11-11'),
    end: new Date('2025-11-12'),
    timed: false,
  })
})
</script>
