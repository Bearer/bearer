import { FetchData, TOAUTH2AuthContext } from '@bearer/functions'

export default class FetchDataRefType extends FetchData implements FetchData<any, any, TOAUTH2AuthContext> {
  action = async (event: any): Promise<any> => {
    return {
      data: []
    }
  }
}
