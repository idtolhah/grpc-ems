require('dotenv').config('../.env.local')
const Sequelize = require('sequelize')

const db = new Sequelize(
  {
    username: process.env.MYSQL_USER,
    password: process.env.MYSQL_PASSWORD,
    database: process.env.MYSQL_DB,
    host: process.env.MYSQL_HOST,
    port: process.env.MYSQL_PORT,
    dialect: 'mysql',
    // operatorsAliases: false,
    // dialectOptions: {
    //   useUTC: false, // for reading from database
    // },
    timezone: '+07:00', // for writing to database
    logging: false,
  }
)

const connectDB = async () => {
  try {
    await db.authenticate()
    console.log('Connection has been established successfully.')
  } catch (error) {
    console.error(`Error: ${error.message}`)
    process.exit(1)
  }
}

module.exports = { db, connectDB }
