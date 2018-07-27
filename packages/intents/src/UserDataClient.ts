import axios, { AxiosInstance } from 'axios'

// class UserDataClient {
//   private client: AxiosInstance
//   constructor(private readonly baseURL: string) {
//     this.client = axios.create({
//       baseURL,
//       timeout: 3000,
//       headers: {
//         Accept: 'application/json',
//         'User-Agent': 'Bearer'
//       }
//     })
//   }

//   retrieveState(referenceId: string) {}
// }

export default (baseURL: string): AxiosInstance =>
  axios.create({
    baseURL,
    timeout: 3000,
    headers: {
      Accept: 'application/json',
      'User-Agent': 'Bearer'
    }
  })
