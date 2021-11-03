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
    return {error};
}
}

export async function getService(serviceId){
    try{
        const {data} = await axiosClient.get(`/${serviceId}`)
        return [null, data];
    }
    catch(error){
        return {error};
    }
}

export async function addService(service){
    try{
        await axiosClient.post(`/servicesadd`,service)
        return [null];
    }
    catch(error){
        return {error};
    }
}

export async function editService(service, serviceId){
    try{
        await axiosClient.put(`/services/edit/${serviceId}`,service)
        return [null];
    }
    catch(error){
        return {error};
    }
}

export async function deleteService(serviceId){
    try{
        await axiosClient.delete(`/services/delete/${serviceId}`)
        return [null];
    }
    catch(error){
        return {error};
    }
}

export async function requestService(serviceId,userId){
    try{
        await axiosClient.post(`/services/request/${serviceId}?uid=${userId}`)
        return [null];
    }
    catch(error){
        return {error};
    }
}



