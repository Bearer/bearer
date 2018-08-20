import axios from 'axios'

export default function(username: string, password: string) {
  const headers = {
    'Accept': 'application/json',
    'User-Agent': 'Bearer'
  }

  return axios.create({
    baseURL: 'https://api.example.com/v1',
    timeout: 5000,
    auth: { username, password },
    headers
  })
}