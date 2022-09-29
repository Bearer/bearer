const orderServiceUrl: string = process.ENV['ORDER_SERVICE_URL'];
let userServiceHost: string = process.ENV.USER_SERVICE_HOST;

const { PRODUCT_SERVICE_URL, AUTH_DOMAIN } = process.ENV;

const accountId = process.ENV.ACCOUNT_ID;
const other = IGNORE_ME_HOST;

const interpolation = `http://${process.ENV.CUSTOMER_HOST}/path/${accountId}`
const concat = 'http://' + process.ENV.CUSTOMER_HOST + "/" + "path/" + accountId

const ignored = <img src="ignored.domain.com"/>
