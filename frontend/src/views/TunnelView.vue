<script setup lang="ts">
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
  // fetchTunnels()
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

const rules = {
  required: (v: any) => !!v || 'Field is required',
  integers: (v: any) => {
    const n = Number(v)
    if (!isNaN(n) && Number.isInteger(n)) {
      return true
    }
    return 'Must be an integer'
  }
}

const editForm = ref(true)
const addForm = ref(false)

onMounted(async () => {
  await fetchTunnels()
  setInterval(fetchTunnels, 10000)
}
)

</script>

<template>
  <v-container class="tunnel-container">
    <v-row v-for="t in tunnels" :key="t.tunnel.id">
      <v-col>
        <v-card @click="enableEditOverlay(t.tunnel)" rounded="xl" class="pa-8 container" variant="flat">
          <v-row class="d-flex align-center">
            <v-col cols="2" class="d-flex justify-center">
              <svg width="16" height="16" xmlns="http://www.w3.org/2000/svg">
                <circle cx="8" cy="8" r="8"
                  :fill="(t.status != undefined) ? ((t.status) ? 'green' : 'red') : 'black'" />
              </svg>
            </v-col>
            <v-col cols="5">
              <h3>{{ t.tunnel.name }}</h3>
            </v-col>
            <v-col cols="5" class="d-flex justify-end">
              <TunnelSwitch v-model="t.tunnel.enabled" @click.stop="" />
            </v-col>
          </v-row>
        </v-card>
      </v-col>
    </v-row>
  </v-container>

  <v-dialog  max-width=500 v-model="editOverlay">
    <template v-slot:default="{ }">
      <v-card rounded="xl">
        <v-card-title>Edit {{ overlayFields.name }}</v-card-title>
        <v-card-text>
          <v-form v-model="editForm">
            <v-row>
              <v-col>
                <v-text-field :rules="[rules.required]" label="Name" v-model="overlayFields.name" variant="outlined"></v-text-field>
              </v-col>
            </v-row>
            <v-row>
              <v-col>
                <v-text-field :rules="[rules.required]" label="Host" v-model="overlayFields.host" variant="outlined"></v-text-field>
              </v-col>
            </v-row>
            <v-row dense>
              <v-col>
                <v-text-field :rules="[rules.required, rules.integers]" label="Local Port" v-model="overlayFields.local_port" variant="outlined"></v-text-field>
              </v-col>
              <v-col>
                <v-text-field :rules="[rules.required, rules.integers]" label="Remote Port" v-model="overlayFields.remote_port" variant="outlined"></v-text-field>
              </v-col>
            </v-row>
            <v-row>
              <v-col>
                <v-text-field :rules="[rules.required]" label="Connection Address" v-model="overlayFields.conn_addr"
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
                <v-btn :color="colors.blue.base" :disabled="!editForm" type="submit" @click="editTunFromOverlay(); editOverlay = false" variant="plain">Save</v-btn>
              </v-col>
            </v-row>
          </v-form>
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
          <v-form v-model="addForm">
            <v-row>
              <v-col>
                <v-text-field :rules="[rules.required]" label="Name" v-model="overlayFields.name" variant="outlined"></v-text-field>
              </v-col>
            </v-row>
            <v-row>
              <v-col>
                <v-text-field :rules="[rules.required]" label="Host" v-model="overlayFields.host" variant="outlined"></v-text-field>
              </v-col>
            </v-row>
            <v-row dense>
              <v-col>
                <v-text-field :rules="[rules.required, rules.integers]" label="Local Port" v-model="overlayFields.local_port" variant="outlined"></v-text-field>
              </v-col>
              <v-col>
                <v-text-field :rules="[rules.required, rules.integers]" label="Remote Port" v-model="overlayFields.remote_port" variant="outlined"></v-text-field>
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
                <v-btn :color="colors.blue.base" :disabled="!addForm" @click="addTunFromOverlay(); addOverlay = false" variant="plain">Save</v-btn>
              </v-col>
            </v-row>
          </v-form>
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
  width: 16px;
  height: 16px;
  border-radius: 30px;
}

.overlay {
  min-width: 300px;
  background-color: var(--color-surface-container);
  color: var(--color-on-surface);
}

.container {
  background-color: var(--color-surface-container);
  color: var(--color-on-surface);
}
</style>