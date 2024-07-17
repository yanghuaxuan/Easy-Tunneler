<script setup lang="ts">
import { mdiWrenchOutline } from "@mdi/js"
import { onMounted, reactive, watch } from "vue";
import { ref } from "vue";
import { mdiPlus } from "@mdi/js"
import colors from 'vuetify/util/colors'
import TunnelSwitch from '@/components/TunnelSwitch.vue'


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

const addr = "http://localhost:4140"

const tunnels = ref<Tunnel_Status[]>([])

const overlayFields = reactive(
  {
    id: "",
    name: "",
    local_port: "",
    host: "",
    remote_port: "",
    conn_addr: "",
    enabled: true,
    autoreboot: false,
  }
)

// const enableEditOverlay = (t: Tunnel) => {
//   overlayFields.id = t.id
//   overlayFields.name = t.name
//   overlayFields.local_port = t.local_port.toString()
//   overlayFields.remote_port = t.remote_port.toString()
//   overlayFields.conn_addr = t.conn_addr
//   overlayFields.autoreboot = t.autoreboot
//   editOverlay.value = true
// }

const enableAddOverlay = () => {
  overlayFields.name = ""
  overlayFields.local_port = ""
  overlayFields.host = "localhost"
  overlayFields.remote_port = ""
  overlayFields.conn_addr = ""
  overlayFields.enabled = true
  overlayFields.autoreboot = true
  // addOverlay.value = true
  addOverlay.value = true
}

const enableEditOverlay = (t: Tunnel) => {
  overlayFields.id = t.id
  overlayFields.name = t.name
  overlayFields.local_port = t.local_port.toString()
  overlayFields.host = t.host
  overlayFields.remote_port = t.remote_port.toString()
  overlayFields.conn_addr = t.conn_addr
  overlayFields.enabled = t.enabled
  overlayFields.autoreboot = t.autoreboot
  // addOverlay.value = true
  editOverlay.value = true
}

const editOverlay = ref(false)
const addOverlay = ref(false)

const editTunFromOverlay = async () => {
  let t: Tunnel = {
    id: overlayFields.id,
    name: overlayFields.name,
    local_port: parseInt(overlayFields.local_port),
    host: overlayFields.host,
    remote_port: parseInt(overlayFields.remote_port),
    conn_addr: overlayFields.conn_addr,
    enabled: overlayFields.enabled,
    autoreboot: overlayFields.autoreboot
  }
  await fetch(`${addr}/api/v1/update_tunnel`, {
    method: "PATCH",
    body: JSON.stringify(t)
  })
  fetchTunnels()
}

const addTunFromOverlay = async () => {
  let t: Tunnel = {
    id: overlayFields.id,
    name: overlayFields.name,
    local_port: parseInt(overlayFields.local_port),
    host: overlayFields.host,
    remote_port: parseInt(overlayFields.remote_port),
    conn_addr: overlayFields.conn_addr,
    enabled: overlayFields.enabled,
    autoreboot: overlayFields.autoreboot
  }
  await fetch(`${addr}/api/v1/add_tunnel`, {
    method: "POST",
    body: JSON.stringify(t)
  })
  fetchTunnels()
}

const deleteTunById = async (tid: string) => {
  await fetch(`${addr}/api/v1/remove_tunnel`, {
    method: "POST",
    body: JSON.stringify({ id: tid })
  })
  await fetchTunnels()
}


const fetchTunnels = async () => {
  await fetch(`${addr}/api/v1/tunnel_status`)
    .then(resp => {
      resp.json().then((j: Resp) => {
        tunnels.value = []
        tunnels.value = j.tunnel_status.sort((a, b) => a.tunnel.name > b.tunnel.name ? 1 : -1)
        for (const t of tunnels.value) {
          watch(t, async (n) => {
            console.log(JSON.stringify(n.tunnel))
            await fetch(`${addr}/api/v1/update_tunnel`, {
              method: "PATCH",
              body: JSON.stringify(n.tunnel)
            })
            fetchTunnels()
          })
        }
      })
    })
    .catch(e => {
      console.error("Cannot fetch tunnel status: " + e)
    })
}

onMounted(async () => {
  await fetchTunnels()
  setInterval(fetchTunnels, 10000)
}
)

</script>

<template>
  <v-container class="tunnel-container">
    <v-row class="d-flex align-center">
      <v-col>
        <h1 class="text-h3">Tunnels</h1>
      </v-col>
      <v-col class="d-flex justify-end">
        <v-btn class="rounded-button" size="x-large" variant="text" :icon="mdiWrenchOutline" />
      </v-col>
    </v-row>
    <v-row v-for="t in tunnels" :key="t.tunnel.id">
      <v-col>
        <v-card @click="enableEditOverlay(t.tunnel)" rounded="xl" class="pa-8 container" variant="flat">
          <v-row class="d-flex align-center">
            <v-col cols="2" class="d-flex justify-center">
              <!-- <div class="dot" :style="{'background-color': (t.status != undefined) ? ((t.status) ? 'green' : 'red') : 'black'}"></div> -->
              <svg width="16" height="16" xmlns="http://www.w3.org/2000/svg">
                <circle cx="8" cy="8" r="8"
                  :fill="(t.status != undefined) ? ((t.status) ? 'green' : 'red') : 'black'" />
              </svg>
            </v-col>
            <v-col>
              <h3>{{ t.tunnel.name }}</h3>
            </v-col>
            <!-- <v-col class="d-flex align-center justify-center">
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
            </v-col> -->
            <!-- <v-col cols="2" class="d-flex flex-row-reverse align-center">
              <v-btn variant="text" @click="enableEditOverlay(t.tunnel)">
                <v-icon :icon="mdiDotsVertical"></v-icon>
              </v-btn>
            </v-col> -->
            <v-col class="d-flex justify-end">
              <!-- <v-switch @click.native.stop v-model="t.tunnel.enabled" color="blue"></v-switch> -->
              <TunnelSwitch v-model="t.tunnel.enabled" @click.stop="" />
            </v-col>
            <!-- <v-col cols="1" class="d-flex justify-end ">
              <v-btn :icon="mdiTrashCanOutline" @click="deleteTun(t.tunnel)" />
            </v-col> -->
          </v-row>
        </v-card>
      </v-col>
    </v-row>
  </v-container>

  <!-- <v-overlay class="justify-center align-center h-screen w-screen" v-model="editOverlay">
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
              <v-text-field label="Connection Address" v-model="overlayFields.conn_addr"
                variant="outlined"></v-text-field>
            </v-col>
          </v-row>
          <v-row>
            <v-col>
              <v-switch label="Autoreboot" color="var(--color-primary-container)" v-model="overlayFields.autoreboot" />
            </v-col>
          </v-row>
          <v-row class="d-flex justify-right">
            <v-col cols="3">
              <v-btn variant="plain" :style="{ color: colors.red.base }">Delete</v-btn>
            </v-col>
            <v-spacer cols />
            <v-col cols="3">
              <v-btn variant="plain">Save</v-btn>
            </v-col>
          </v-row>
        </v-container>
      </v-card>
    </v-overlay> -->
  <v-dialog v-model="editOverlay">
    <template v-slot:default="{ }">
      <v-card rounded="xl">
        <v-card-title>Edit {{ overlayFields.name }}</v-card-title>
        <v-card-text>
          <v-row>
            <v-col>
              <v-text-field label="Name" v-model="overlayFields.name" variant="outlined"></v-text-field>
            </v-col>
          </v-row>
          <v-row>
            <v-col>
              <v-text-field label="Host" v-model="overlayFields.host" variant="outlined"></v-text-field>
            </v-col>
          </v-row>
          <v-row dense>
            <v-col>
              <v-text-field label="Local Port" v-model="overlayFields.local_port" variant="outlined"></v-text-field>
            </v-col>
            <v-col>
              <v-text-field label="Remote Port" v-model="overlayFields.remote_port" variant="outlined"></v-text-field>
            </v-col>
          </v-row>
          <v-row>
            <v-col>
              <v-text-field label="Connection Address" v-model="overlayFields.conn_addr"
                variant="outlined"></v-text-field>
            </v-col>
          </v-row>
          <v-row>
            <v-col>
              <v-switch label="Autoreboot" color="var(--color-primary-container)" v-model="overlayFields.autoreboot" />
            </v-col>
          </v-row>
          <v-row class="d-flex justify-right">
            <v-col cols="3">
              <v-btn @click="deleteTunById(overlayFields.id); editOverlay = false" variant="plain"
                :style="{ color: colors.red.base }">Delete</v-btn>
            </v-col>
            <v-spacer cols />
            <v-col cols="3" class="d-flex justify-end">
              <v-btn @click="editTunFromOverlay(); editOverlay = false" variant="plain">Save</v-btn>
            </v-col>
          </v-row>
        </v-card-text>
      </v-card>
    </template>
  </v-dialog>
  <v-fab size="72" color="var(--color-primary-container)" @click="enableAddOverlay()" :icon="mdiPlus" rounded="xl" app
    location="bottom end" class="mr-6 mb-12"></v-fab>
  <v-dialog max-width=500 v-model="addOverlay">
    <template v-slot:default="{ }">
      <v-card rounded="xl">
        <v-card-title>Add tunnel</v-card-title>
        <v-card-text>
          <v-row>
            <v-col>
              <v-text-field label="Name" v-model="overlayFields.name" variant="outlined"></v-text-field>
            </v-col>
          </v-row>
          <v-row>
            <v-col>
              <v-text-field label="Host" v-model="overlayFields.host" variant="outlined"></v-text-field>
            </v-col>
          </v-row>
          <v-row dense>
            <v-col>
              <v-text-field label="Local Port" v-model="overlayFields.local_port" variant="outlined"></v-text-field>
            </v-col>
            <v-col>
              <v-text-field label="Remote Port" v-model="overlayFields.remote_port" variant="outlined"></v-text-field>
            </v-col>
          </v-row>
          <v-row>
            <v-col>
              <v-text-field placeholder="john@example.com" label="Connection Address" v-model="overlayFields.conn_addr"
                variant="outlined"></v-text-field>
            </v-col>
          </v-row>
          <v-row>
            <v-col>
              <v-switch label="Autoreboot" v-model="overlayFields.autoreboot"
                color="var(--color-primary-container)">Autoreboot</v-switch>
            </v-col>
          </v-row>
          <v-row>
            <v-col class="d-flex justify-end">
              <v-btn @click="addTunFromOverlay(); addOverlay = false" variant="plain">Save</v-btn>
            </v-col>
          </v-row>
        </v-card-text>
      </v-card>
    </template>
  </v-dialog>
</template>

<style scoped>
.v-card-title {
  font-size: 1.6em;
  font-weight: bold;
}

.tunnel-container {
  max-width: 700px;
}

.dot {
  display: block;
  /* background-color: rgb(2, 179, 2); */
  width: 16px;
  height: 16px;
  border-radius: 30px;
}

.overlay {
  /* width: 30vw; */
  /* max-width: 400px; */
  min-width: 300px;
  background-color: var(--color-surface-container);
  color: var(--color-on-surface);
}

.container {
  background-color: var(--color-surface-container);
  color: var(--color-on-surface);
}
</style>