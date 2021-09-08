function baseUrl(){
    return "";
}

function parseUrlParameters(location){
    return new Map(location.search.slice(1).split('&').map(kv => kv.split('=')));
}

function isNullOrUndefined(element){
    return element == null || element == undefined
}

async function $ajax(option){
    const obj = {
        type: "post",
        contentType: "application/json",
        dataType: "json",
        data: `{}`,
        async: false,
    }
    Object.assign(obj , option)
    return new Promise ((resolve, rejact) => {
        $.ajax({
            ...obj,
            success(data){
                resolve(data)
            },
            error(err){
                rejact(err)
            }
        })
    })
}

function transactionHtml(item){
    let {transactionInputs,transactionOutputs} = item;

    let left = '';
    if(!isNullOrUndefined(transactionInputs)){
        transactionInputs.forEach(item1 => {
            left += `<div>付：<span><a title="地址，点击查看地址详情。" target="_blank" href="./address.html?address=${item1.address}">${item1.address}</a></span>&nbsp;<a title="交易输出，点击查看交易输出详情。" target="_blank" href="./transactionoutput.html?transactionHash=${item1.transactionHash}&transactionOutputIndex=${item1.transactionOutputIndex}"><i class="glyphicon-euro"></i></a>&nbsp;<span>${item1.value}</span></div>`
        })
    }

    let right = ''
    if(!isNullOrUndefined(transactionOutputs)){
        transactionOutputs.forEach(item1 => {
            right += `<div style="display:flex">收：<span><a title="地址，点击查看地址详情。" target="_blank" href="./address.html?address=${item1.address}">${item1.address}</a></span>&nbsp;<a title="交易输出，点击查看交易输出详情。" target="_blank" href="./transactionoutput.html?transactionHash=${item1.transactionHash}&transactionOutputIndex=${item1.transactionOutputIndex}"><i class="glyphicon-euro"></i></a>&nbsp;<span>${item1.value}</span></div>`
        });
    }

    let transactionHtml = `
        <div style="font-size: 14px; line-height: 40px; margin-top: 10px;">
            <!-- 开头 -->
            <div style="display: flex; background-color: #f5f5f5; flex-wrap: wrap; border-top: 1px solid #ddd; border-bottom: 1px solid #ddd;">
                <div style="min-width: 40%;"><a title="交易哈希，点击查看交易详情。" target="_blank" href="./transaction.html?transactionHash=${item.transactionHash}">${item.transactionHash}</a></div>
                <div style="min-width: 20%; text-align: center;">${item.transactionType}</div>
                <div style="min-width: 20%; text-align: center;">${item.transactionFee}</div>
                <div style="min-width: 20%; text-align: center;">${item.blockTime}</div>
            </div>
            <div style="display: flex; flex-wrap:wrap;">
                <div style="width:auto; min-width:50%;">${left}</div>
                <div style="width:auto; min-width:40%;">${right}</div>
            </div>
        </div>
    `;
    return transactionHtml;
}

function unconfirmedTransactionHtml(item){
    let {transactionInputs,transactionOutputs} = item;

    let left = '';
    if(!isNullOrUndefined(transactionInputs)){
        transactionInputs.forEach(item1 => {
            left += `<div>付：<span><a title="地址，点击查看地址详情。" target="_blank" href="./address.html?address=${item1.address}">${item1.address}</a></span>&nbsp;<a title="交易输出，点击查看交易输出详情。" target="_blank" href="./transactionoutput.html?transactionHash=${item1.transactionHash}&transactionOutputIndex=${item1.transactionOutputIndex}"><i class="glyphicon-euro"></i></a>&nbsp;<span>${item1.value}</span></div>`
        })
    }

    let right = ''
    if(!isNullOrUndefined(transactionOutputs)){
        transactionOutputs.forEach(item1 => {
            right += `<div style="display:flex">收：<span><a title="地址，点击查看地址详情。" target="_blank" href="./address.html?address=${item1.address}">${item1.address}</a></span>&nbsp;<a title="交易输出，点击查看交易输出详情。" target="_blank" href="./transactionoutput.html?transactionHash=${item1.transactionHash}&transactionOutputIndex=${item1.transactionOutputIndex}"><i class="glyphicon-euro"></i></a>&nbsp;<span>${item1.value}</span></div>`
        });
    }

    let transactionHtml = `
        <div style="font-size: 14px; line-height: 40px; margin-top: 10px;">
            <!-- 开头 -->
            <div style="display: flex; background-color: #f5f5f5; flex-wrap: wrap; border-top: 1px solid #ddd; border-bottom: 1px solid #ddd;">
                <div style="min-width: 40%;"><a title="交易哈希，点击查看交易详情。" target="_blank" href="./transaction.html?transactionHash=${item.transactionHash}">${item.transactionHash}</a></div>
            </div>
            <div style="display: flex; flex-wrap:wrap;">
                <div style="width:auto; min-width:50%;">${left}</div>
                <div style="width:auto; min-width:40%;">${right}</div>
            </div>
        </div>
    `;
    return transactionHtml;
}

function localLanguage(element){
    if(element == 'REQUEST_PARAM_ILLEGAL'){
        return '请求参数非法';
    }
    if(element == 'SERVICE_UNAVAILABLE'){
        return '服务不可用';
    }

    if(element == 'BUILD_TRANSACTION_SUCCESS'){
       return '构建交易成功';
    }
    if(element == 'PAYEE_CAN_NOT_EMPTY'){
        return '收款方不能为空';
    }
    if(element == 'PAYEE_ADDRESS_CAN_NOT_EMPTY'){
        return '收款方地址不能为空';
    }
    if(element == 'PAYEE_VALUE_CAN_NOT_LESS_EQUAL_THAN_ZERO'){
        return '收款金额不能小于等于零';
    }
    if(element == 'NOT_ENOUGH_MONEY_TO_PAY'){
       return '没有足够的钱进行支付';
    }
    return '本地语言翻译异常,请联系管理员';
}