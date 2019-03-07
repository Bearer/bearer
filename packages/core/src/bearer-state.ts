import * as requests from './requests'

export function storeSetup(payload: any) {
  return postSetup({ ...payload, ReadAllowed: false })
}

export function storeSecret(referenceId: string, payload: any) {
  return putItem(referenceId, { ...payload, ReadAllowed: false })
}

export function storeData(referenceId: string, payload: any) {
  return putItem(referenceId, { ...payload, ReadAllowed: true })
}

export function getData(referenceId: string) {
  const request = requests.functionRequest({
    functionName: referenceId,
    integrationId: 'items',
    setupId: 'TODO'
  })
  return request({}, {})
}

export function removeData(referenceId: string) {
  const request = requests.functionRequest({
    functionName: referenceId,
    integrationId: 'items',
    setupId: 'TODO'
  })
  return request(
    {},
    {
      method: 'DELETE'
    }
  )
}

function postSetup(payload: any) {
  const request = requests.itemRequest()
  return request(
    {},
    {
      method: 'POST',
      body: JSON.stringify(payload)
    }
  )
}

function putItem(referenceId: string, payload: any) {
  const request = requests.functionRequest({
    functionName: referenceId,
    integrationId: 'items',
    setupId: 'TODO'
  })
  return request(
    {},
    {
      method: 'PUT',
      body: JSON.stringify(payload)
    }
  )
}
