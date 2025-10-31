export default {
    user: {
      authenticated: false,
      name: '',
    },
  
    logout () {
      localStorage.removeItem('token')
      this.user.authenticated = false
      this.user.name = ''
    },
  }