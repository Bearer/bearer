import axios from 'axios'

export const CLIENT = axios.create({
  baseURL: 'https://swapi.co/api/',
  timeout: 5000,
  headers: {
    Accept: 'application/json',
    'User-Agent': 'Bearer'
  }
})

export function headersFor(token) {
  return {
    Authorization: `token ${token}`
  }
}
