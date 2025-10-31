<template>
  <v-container fluid>
    <v-responsive>
      <v-row>
        <v-col cols="auto">
          <v-card>
            <v-card-text class="pa-0">

              <v-tabs v-model="selectedTab" @update:model-value="tabChanged" class="rounded">
                <v-tab class="ml-2" value="queries">Queries</v-tab>
                <v-tab value="tablesize">Table size</v-tab>
                <v-tab value="settings">Settings</v-tab>
                <v-tab value="unusedidxs">Unused indexes</v-tab>
                <v-tab value="bloat">Bloat</v-tab>
              </v-tabs>

              <v-window v-model="selectedTab">
                
                <v-window-item value="queries">
                  <TechPgQueryTable />
                </v-window-item>

                <v-window-item value="tablesize">
                  <TechPgTableSizeTable />
                </v-window-item>

                <v-window-item value="settings">
                  <TechPgSettingTable />
                </v-window-item>

                <v-window-item value="unusedidxs">
                  <TechPgUnusedIdxTable />
                </v-window-item>

                <v-window-item value="bloat">
                  <TechPgBloatTable />
                </v-window-item>

              </v-window>
              
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>
    </v-responsive>
  </v-container>
</template>

<script lang="ts" setup>
import { ref, onBeforeMount } from 'vue'

const lsKey = 'tech_database'
const selectedTab = ref('queries')

function tabChanged() {
  let lsObj = {
    'selectedTab': selectedTab.value,
  }
  localStorage.setItem(lsKey, JSON.stringify(lsObj))
}

onBeforeMount(() => {
  var lsJSON = localStorage.getItem(lsKey)
  if (!lsJSON) {
    return
  }

  let lsObj = JSON.parse(lsJSON)
  if (lsObj['selectedTab']) { selectedTab.value = lsObj['selectedTab'] }
})
</script>
