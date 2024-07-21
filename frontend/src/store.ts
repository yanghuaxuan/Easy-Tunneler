import { reactive } from 'vue'

export const store = reactive({
    autorefresh_interval: parseInt(localStorage.getItem("autorefresh_interval") ?? "15"),
})