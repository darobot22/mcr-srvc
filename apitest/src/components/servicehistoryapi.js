import axios from 'axios';
const axiosClient = axios.create({
    baseURL: 'http://localhost:1236'
});

export async function getAllHistory(){
    try{
        const {data} = await axiosClient.get(`/history`)
        return [null, data];
    }
    catch(error){
        return {error};
    }
}

export async function getHistory(userId){
    try{
        const {data} = await axiosClient.get(`/history/${userId}`)
        return [null, data];
    }
    catch(error){
        return {error};
    }
}

export async function editHistory(service, historyId){
    try{
        await axiosClient.put(`/history/edit/${historyId}`,service)
        return [null];
    }
    catch(error){
        return {error};
    }
}