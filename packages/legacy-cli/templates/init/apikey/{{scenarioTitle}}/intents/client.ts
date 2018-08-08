import axios from 'axios'

export default function(token: string) {
  const encodedAuth = Buffer.from(`anything:${token}`).toString('base64')
  const headers = {
    'Accept': 'application/json',
    'User-Agent': 'Bearer',
    'Authorization': `Basic ${encodedAuth}`
  }

  return axios.create({
    baseURL: 'https://api.example.com/v1',
    timeout: 5000,
    headers
  })
}