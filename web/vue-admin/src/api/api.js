import axios from 'axios'

let base = ''
console.log(process.env)
if (process.env.NODE_ENV === 'development') {
  base = 'http://127.0.0.1:8080/api/v1'
} else {
  base = 'api/v1'
}

export const requestLogin = params => { return axios.post(`${base}/login`, params).then(res => res.data) }

export const getUserList = params => { return axios.get(`${base}/user/list`, params) }

export const getUserListPage = params => { return axios.get(`${base}/user/listpage`, { params: params }) }

export const removeUser = params => { return axios.post(`${base}/user/remove`, params) }

export const batchRemoveUser = params => { return axios.post(`${base}/user/batchremove`, params) }

export const editUser = params => { return axios.post(`${base}/user/edit`, params).then(res => res.data) }

export const addUser = params => { return axios.post(`${base}/user/add`, params).then(res => res.data) }

export const getAllUsers = params => { return axios.get(`${base}/user/alluid`, params) }

export const getUserInfo = params => { return axios.get(`${base}/user/info`, { params: params }) }

export const getRotaListPage = params => { return axios.get(`${base}/rota/listpage`, { params: params }) }

export const removeRota = params => { return axios.post(`${base}/rota/remove`, params) }

export const batchRemoveRota = params => { return axios.post(`${base}/rota/batchremove`, params) }

export const editRota = params => { return axios.post(`${base}/rota/edit`, params).then(res => res.data) }

export const addRota = params => { return axios.post(`${base}/rota/add`, params).then(res => res.data) }

export const getReceiverInfoListPage = params => { return axios.get(`${base}/receiverInfo/listpage`, { params: params }) }

export const removeReceiverInfo = params => { return axios.post(`${base}/receiverInfo/remove`, params) }

export const batchRemoveReceiverInfo = params => { return axios.post(`${base}/receiverInfo/batchremove`, params) }

export const editReceiverInfo = params => { return axios.post(`${base}/receiverInfo/edit`, params).then(res => res.data) }

export const addReceiverInfo = params => { return axios.post(`${base}/receiverInfo/add`, params).then(res => res.data) }
