const express = require("express")
const helmet = require("helmet")

const app = express()
app.use(helmet())
// app.use(helmet.contentSecurityPolicy());
// app.use(helmet.crossOriginEmbedderPolicy());
// app.use(helmet.crossOriginOpenerPolicy());
// app.use(helmet.crossOriginResourcePolicy());
// app.use(helmet.dnsPrefetchControl());
// app.use(helmet.frameguard());
// app.use(helmet.hsts());
// app.use(helmet.ieNoOpen());
// app.use(helmet.noSniff());
// app.use(helmet.originAgentCluster());
// app.use(helmet.permittedCrossDomainPolicies());
// app.use(helmet.referrerPolicy());
// app.use(helmet.xssFilter());
// app.use(helmet.hidePoweredBy());

app.disable("x-powered-by")
