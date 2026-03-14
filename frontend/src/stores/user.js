import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useUserStore = defineStore('user', () => {
  const user = ref(null)

  const isLoggedIn = computed(() => !!user.value)
  const isAdmin = computed(() => user.value?.is_admin === true)

  function setUser(data) {
    user.value = data
  }

  function logout() {
    user.value = null
  }

  return { user, isLoggedIn, isAdmin, setUser, logout }
})
