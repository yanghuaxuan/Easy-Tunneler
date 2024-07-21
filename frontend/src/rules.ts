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

export default rules