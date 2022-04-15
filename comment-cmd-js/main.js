require('dotenv').config({ path: './.env.local' })
var PROTO_PATH = __dirname + '/commentcmdpb/commentcmd.proto'

var grpc = require('@grpc/grpc-js')
var protoLoader = require('@grpc/proto-loader')
const moment = require('moment')
const { connectDB, db } = require('./db/mysql')
var packageDefinition = protoLoader.loadSync(
    PROTO_PATH,
    {keepCase: true,
     longs: String,
     enums: String,
     defaults: true,
     oneofs: true
    })
var comment_proto = grpc.loadPackageDefinition(packageDefinition).comment

async function createPackingComment(call, callback) {
    console.log('Called CreatePackingComment: ')
    console.log('Call: ' + JSON.stringify(call.request))

    try {
        const req = call.request
        const [results, metadata] = await db.query(`
            UPDATE equipment_checkings SET mr_id = '${req.mr_id}', mr_comment = '${req.mr_comment}', mr_created_at = '${moment().format('YYYY-MM-DD hh:mm:ss')}' WHERE id = ${req.equipment_checking_id}
        `);

        console.log(JSON.stringify(results))
        callback(null, {results})
    } catch (error) {
        console.log(error)
    }
}

function main() {
    console.log('Comment Command Service')
    
    var server = new grpc.Server();
    const addr = `${process.env.GRPC_SERVICE_HOST}:${process.env.GRPC_SERVICE_PORT}`
    server.addService(comment_proto.CommentCmdService.service, {createPackingComment})
    server.bindAsync(addr, grpc.ServerCredentials.createInsecure(), () => {
        server.start() 
        console.log(`Server started at ${addr}`)
    });

    connectDB()
}

main()
