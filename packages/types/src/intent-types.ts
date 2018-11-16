export enum IntentType {
  /**
   * Simple intent (function) to fetch data (get/post/put) from
   * You can store data within the intent itself or query external APIs
   */
  FetchData = 'FetchData',
  /**
   * Let you store data easily data to the Bearer database
   */
  SaveState = 'SaveState',
  /**
   * Let you store retrieve easily data from the Bearer database
   */
  RetrieveState = 'RetrieveState'
}

export default IntentType
