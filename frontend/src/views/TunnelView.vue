<script setup lang="ts">
import { mdiPlus } from "@mdi/js"
import { onMounted, reactive } from "vue";
import { ref } from "vue";
import { mdiArrowRightBold, mdiTrashCanOutline } from "@mdi/js"
import colors from 'vuetify/util/colors'


interface Resp {
  tunnel_status: Tunnel_Status[]
}

interface Tunnel_Status {
  tunnel: Tunnel
  status?: number
}

interface Tunnel {
  conn_addr: string
  enabled: boolean
  host: string
  id: string
  local_port: number
  name: string
  remote_port: number
  autoreboot: boolean 
}

const tunnels = ref<Tunnel_Status[]>([])

const overlayFields = reactive(
  {
    id: "",
    name: "",
    local_port: "",
    remote_port: "",
    conn_addr: "",
    autoreboot: false
  }
)

const enableEditOverlay = (t: Tunnel) => {
  overlayFields.id = t.id
  overlayFields.name = t.name
  overlayFields.local_port = t.local_port.toString()
  overlayFields.remote_port = t.remote_port.toString()
  overlayFields.conn_addr = t.conn_addr
  overlayFields.autoreboot = t.autoreboot
  editOverlay.value = true
}

const enableAddOverlay = () => {
  overlayFields.id = ""
  overlayFields.name = ""
  overlayFields.local_port = ""
  overlayFields.remote_port = ""
  overlayFields.conn_addr = ""
  overlayFields.autoreboot = true
  addOverlay.value = true
}

const editOverlay = ref(false)
const addOverlay = ref(false)

const deleteTun = async (t: Tunnel) => {
  await fetch("http://localhost:4140/api/v1/remove_tunnel", {
    method: "POST",
    body: JSON.stringify({id: t.id})
  })
  await fetchTunnels()
}

const fetchTunnels = async () => {
  console.log("fetching")
  await fetch("http://localhost:4140/api/v1/tunnel_status")
    .then(resp => {
      resp.json().then((j: Resp) => {
        tunnels.value = j.tunnel_status.sort((a,b) =>  a.tunnel.name > b.tunnel.name ? 1 : -1)
      })
    })
    .catch(e => {
      console.error("Cannot fetch tunnel status: " + e)
    })
}

onMounted(async () =>  {
    await fetchTunnels()
    setInterval(fetchTunnels, 10000)
  }
)

</script>

<template>
  <main>
    <v-row>
      <v-col>
        <v-card class="pa-3" variant="outlined" rounded="lg">
          <v-row>
            <v-col>
              <h1>Tunnels</h1>
            </v-col>
            <v-col class="d-flex justify-end align-center">
              <v-btn variant="outlined" rounded @click="enableAddOverlay()"><v-icon size="30" :icon="mdiPlus"></v-icon></v-btn>
            </v-col>
          </v-row>
        </v-card>
      </v-col>
    </v-row>
    <v-row>
      <v-col class="d-flex justify-center">
        <v-card width="95%" class="pa-3" variant="outlined" rounded="xl">
          <v-row v-for="t in tunnels" :key="t.tunnel.id">
            <v-col cols="1" class="d-flex align-center justify-center">
              <div class="dot" :style="{'background-color': (t.status != undefined) ? ((t.status) ? 'green' : 'red') : 'black'}"></div>
            </v-col>
            <v-col class="d-flex align-center">
              <h3>{{ t.tunnel.name }}</h3>
            </v-col>
            <v-col class="d-flex align-center justify-center">
              <v-row>
                <v-col>
                  <p>{{ t.tunnel.local_port }} <v-icon :icon="mdiArrowRightBold" /> {{ t.tunnel.remote_port }}</p>
                  <p>{{ t.tunnel.conn_addr }}</p>
                </v-col>
              </v-row>
              <v-row>
                <v-col>
                </v-col>
              </v-row>
            </v-col>
            <!-- <v-col cols="2" class="d-flex flex-row-reverse align-center">
              <v-btn variant="text" @click="enableEditOverlay(t.tunnel)">
                <v-icon :icon="mdiDotsVertical"></v-icon>
              </v-btn>
            </v-col> -->
            <v-col cols="2" class="d-flex flex-row-reverse align-center">
              <v-btn variant="text" @click="deleteTun(t.tunnel)" >
                <v-icon size="25" :icon="mdiTrashCanOutline"></v-icon>
              </v-btn>
            </v-col>
          </v-row>
        </v-card>
      </v-col>
    </v-row>
    <v-overlay class="justify-center align-center h-screen w-screen" v-model="editOverlay">
      <v-card class="pa-3 overlay" rounded="xl">
        <v-container>
          <v-row>
            <v-col>
              <h1>Edit {{ overlayFields.name }}</h1>
            </v-col>
          </v-row>
          <v-row>
            <v-col>
              <v-text-field label="Local Port" v-model="overlayFields.local_port" variant="outlined"></v-text-field>
            </v-col>
            <v-col>
              <v-text-field label="Remote Port" v-model="overlayFields.remote_port" variant="outlined"></v-text-field>
            </v-col>
          </v-row>
          <v-row>
            <v-col>
              <v-text-field label="Connection Address" v-model="overlayFields.conn_addr" variant="outlined"></v-text-field>
            </v-col>
          </v-row>
          <v-row>
            <v-col>
              <v-switch label="Autoreboot" color="var(--color-primary-container)" v-model="overlayFields.autoreboot" />
            </v-col>
          </v-row>
          <v-row class="d-flex justify-right">
            <v-col cols="3">
              <v-btn variant="plain" :style="{ color: colors.red.base}">Delete</v-btn>
            </v-col>
            <v-spacer cols />
            <v-col cols="3">
              <v-btn variant="plain">Save</v-btn>
            </v-col>
          </v-row>
        </v-container>
      </v-card>
    </v-overlay>
    <v-overlay class="justify-center align-center h-screen w-screen" v-model="addOverlay">
      <v-card class="pa-3 overlay" rounded="xl">
        <v-container>
          <v-row>
            <v-col>
              <h1>Add tunnel</h1>
            </v-col>
          </v-row>
          <v-row>
            <v-col>
              <v-text-field label="Local Port" v-model="overlayFields.local_port" variant="outlined"></v-text-field>
            </v-col>
            <v-col>
              <v-text-field label="Remote Port" v-model="overlayFields.remote_port" variant="outlined"></v-text-field>
            </v-col>
          </v-row>
          <v-row>
            <v-col>
              <v-text-field label="Connection Address" v-model="overlayFields.conn_addr" variant="outlined"></v-text-field>
            </v-col>
          </v-row>
          <v-row>
            <v-col>
              <v-switch label="Autoreboot" v-model="overlayFields.autoreboot" color="var(--color-primary-container)">Autoreboot</v-switch>
            </v-col>
          </v-row>
          <v-row class="d-flex justify-right">
            <v-col cols="3">
              <v-btn variant="plain" :style="{ color: colors.red.base}">Delete</v-btn>
            </v-col>
            <v-spacer cols />
            <v-col cols="3">
              <v-btn variant="plain">Save</v-btn>
            </v-col>
          </v-row>
        </v-container>
      </v-card>
    </v-overlay>
  </main>
</template>

<style scoped>
.dot {
  display: block;
  /* background-color: rgb(2, 179, 2); */
  width: 16px;
  height: 16px;
  border-radius: 30px;
}
.overlay {
  width: 30vw;
  /* max-width: 400px; */
  min-width: 300px;
  background-color: var(--color-surface-container);
  color: var(--color-on-surface);
}
</style>