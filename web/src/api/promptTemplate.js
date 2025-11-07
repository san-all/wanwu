import service from "@/utils/request"
const BASE_URL = '/user/api/v1'
//获取自定义propmpt详情
export const getPromptTemplateDetail = (data)=>{
    return service({
        url: `${BASE_URL}/prompt/custom`,
        method: 'get',
        params: data
    })
}
//获取自定义prompt列表
export const getPromptTemplateList= (data)=>{
    return service({
        url: `${BASE_URL}/prompt/custom/list`,
        method: 'get',
        params: data
    })
}

//获取内置prompt列表
export const getPromptBuiltInList= (data)=>{
    return service({
        url: `${BASE_URL}/prompt/template/list`,
        method: 'get',
        params: data
    })
}
//获取内置prompt详情
export const getPromptBuiltInDetail= (data)=>{
    return service({
        url: `${BASE_URL}/prompt/template/detail`,
        method: 'get',
        params: data
    })
}

