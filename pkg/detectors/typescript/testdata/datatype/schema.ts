// supported
interface User {
    name: string;
    id: number;
    bol: boolean | null
}

// note this type and above type User get merged into 1 we log as seperate
interface User {
    lastname: string
}

// supported simple type literals
export type UserID = "testuser"


// testuser is kept as seperate schema (extends for classes and interfaces is not supported)
interface TestUser extends User {
    test_id: string
}

//supported
export type Cat = { breed: string; yearOfBirth: number };

// complex type literlas datatypes are not supported
export type Complex = {
  name: "namevalue"
  surname: "surname"
}

// fully supported with prettyName as string
export type CloudinaryImageFormat = {
  prettyName?: string | null;
  roomInfo: {
      room_id: number;
      short_id: number;
  };
};
  
// supported as object properties address and address2 of type object
interface Addressable {
  address: string[];
  address2: Array<string> ;
  item: {
    id: string;
    channel: string;
    zip: {
      code: string;
    };
  };}


/// supported
export class CloudinaryAdapter {
  // cloudname: string gets joined to cloudfinaryAdapter
  cloudName: string;

  // supported as constructor method of cloudinaryAdapter
  constructor({
    cloudName,
    apiKey,
    apiSecret,
    folder,
  }: {
    cloudName: string;
    apiKey: string;
    apiSecret: string;
    folder?: string;
  }) {
  }

  // supported
  save(file: string) {
  }
}

// supported as Base with property name of type string
class Base {
  name = "base";
  constructor() {
    console.log("My name is " + this.name);
  }
}
   
// extends parts are not supported
// supported as object Dervied with property name of type string
class Derived extends Base {
  name = "derived";
}

// supported with dataypes unkown
class Box<Type> {
  contents: Type;
  // optional parameters are also exported
  constructor(value: Type, optional?: string) {
    this.contents = value;
  }

  save(something: {parent:{nested:string}}) {}
}

// this one is parsed as class instead of class declaration
export default class EventManager {
  private RNOneSignal: NativeModule;
  private oneSignalEventEmitter: NativeEventEmitter;
  private eventHandlerMap: Map<string, (event: any) => void>;
  private eventHandlerArrayMap: Map<string, Array<(event: any) => void>>;
  private listeners: { [key: string]: EmitterSubscription }

  // // empty object
  static getTags(handler: (tags: { [key: string]: string }) => void): void {
  }

  // sms auth code not identified as string
  static setSMSNumber(smsNumber: string, smsAuthCode?: string | null, handler?: Function): void {
  }
}
   