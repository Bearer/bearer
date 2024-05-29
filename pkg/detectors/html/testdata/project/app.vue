<template></template>

<script>
export default {
  data() {
    return {
      posts: [],
    }
  },

  methods: {
    async getData() {
      try {
        const response = await this.$http.get(
          "http://jsonplaceholder.typicode.com/posts"
        )
        // JSON responses are automatically parsed.
        this.posts = response.data
      } catch (error) {
        console.log(error)
      }
    },
  },
}
</script>
