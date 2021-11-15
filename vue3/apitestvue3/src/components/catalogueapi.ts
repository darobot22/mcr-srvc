import axios from 'axios';

const axiosClient = axios.create({
    baseURL: 'http://localhost:1235'
});

export async function getAllServices(){
try{
    const {data} = await axiosClient.get(`/services`)
    return [null, data];
}
catch(error){
    return [error];
}
}

export async function getService(serviceId: number){
    try{
        const {data} = await axiosClient.get(`/${serviceId}`)
        return [null, data];
    }
    catch(error){
        return [error];
    }
}

export async function addService(service: JSON){
    try{
        await axiosClient.post(`/servicesadd`,service)
        return [null];
    }
    catch(error){
        return [error];
    }
}

export async function editService(service: JSON, serviceId: number){
    try{
        await axiosClient.put(`/services/edit/${serviceId}`,service)
        return [null];
    }
    catch(error){
        return [error];
    }
}

export async function deleteService(serviceId: number){
    try{
        await axiosClient.delete(`/services/delete/${serviceId}`)
        return [null];
    }
    catch(error){
        return [error];
    }
}

export async function requestService(serviceId: number,userId: number){
    try{
        await axiosClient.post(`/services/request/${serviceId}?uid=${userId}`)
        return [null];
    }
    catch(error){
        return [error];
    }
}



