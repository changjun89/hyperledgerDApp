const express = require('express');
const bodyPaser = require('body-parser');
const fs = require('fs');
const app = express();

app.use(bodyPaser.urlencoded({extended:false}));
app.set('view engine','jade');
app.set('views','./views');
app.locals.pretty = true;

// Hyperledger Bridge
const { FileSystemWallet, Gateway } = require('fabric-network');
const path = require('path');
const ccpPath = path.resolve(__dirname, '..', 'network' ,'connection.json');
const ccpJSON = fs.readFileSync(ccpPath, 'utf8');
const ccp = JSON.parse(ccpJSON);

async function getContract(){
    
    const walletPath = path.join(process.cwd(), 'wallet');
    const wallet = new FileSystemWallet(walletPath);

    const userExists = await wallet.exists('user1');
    if (!userExists) {
        console.log('An identity for the user "user1" does not exist in the wallet');
        console.log('Run the registerUser.js application before retrying');
        return;
    }
    const gateway = new Gateway();
    await gateway.connect(ccp, { wallet, identity: 'user1', discovery: { enabled: false } });
    const network = await gateway.getNetwork('mychannel');
    const contract = network.getContract('auction');
    return contract;
}

app.get(['/','/goods'],async (req,res) => {
    const contract = await getContract();
    const result = await contract.evaluateTransaction('queryAllGoods');
    const goods = JSON.parse(result);
    
    res.render('view',{'goods':goods});
});


app.get('/goods/:id',async (req,res)=>{
    const id = req.params.id;

    const contract = await getContract();
    const result = await contract.evaluateTransaction('queryGoods',id);
    const goods = JSON.parse(result);

    res.render('detail',{'goods':goods,'bidInfos':goods.bidInfos});
});


app.post('/bidInfo/:id',async (req,res) => {
    const id = req.params.id;
    const userId = req.body.id;
    const price = req.body.price;
    
    const contract = await getContract();
    await contract.submitTransaction('updateBidInfo',id,userId,price);

    return res.redirect("/goods/"+id);
});

app.post('/goods/:id',async (req,res) => {
    const id = req.params.id;        
    
    const contract = await getContract();
    const result = await contract.evaluateTransaction('queryGoods',id);
    const goods = JSON.parse(result);
    const bidInfos = goods.bidInfos;

    var max = bidInfos.reduce( function (previous, current) { 
        return previous.price > current.price ? previous:current;
    });

    await contract.submitTransaction('updateWinUser',id,max.userId,max.price);
    
    return res.redirect("/goods/"+id);
});

app.post('/goods',async (req,res) => {
    const id = req.body.id;
    const name = req.body.name;

    const contract = await getContract();
    await contract.submitTransaction('addGoods',id,name);
    
    return res.redirect("/goods");
});

app.get('/form',async (req,res) => {
    
    res.render('form',{});
});


app.listen(8080,() => {
    console.log('conneted, 8080 port');
});