export enum IntentType {
  /**
   * Simple intent (function) to fetch data (get/post/put) from
   * You can store data within the intent itself or query external APIs
   */
  FetchData = 'FetchData',
  /**
   * Lets you store data easily data to the Bearer database
   */
  SaveState = 'SaveState',
  /**
   * Lets you retrieve data from the Bearer database
   */
  RetrieveState = 'RetrieveState'
}

export default IntentType
