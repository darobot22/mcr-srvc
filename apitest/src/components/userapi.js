import axios from 'axios';

const axiosClient = axios.create({
    baseURL: 'http://localhost:1234',
});

export async function getAllUsers(){
    try{
        const {data} = await axiosClient.get(`/users`)
        return [null, data];
    }
    catch(error){
        return {error};
    }
}

export async function getUser(userId){
    try{
        const {data} = await axiosClient.get(`users/${userId}`)
        return [data];
    }
    catch(error){
        return {error};
    }
}

export async function addUser(user){
    try{
        await axiosClient.post(`/usersadd`,user)
        return [null];
    }
    catch(error){
        return {error};
    }
}

export async function editUser(user, userId){
    try{
        await axiosClient.put(`/users/edit/${userId}`,user)
        return [null];
    }
    catch(error){
        return {error};
    }
}

export async function deleteUser(userId){
    try{
        await axiosClient.delete(`/users/delete/${userId}`)
        return [null];
    }
    catch(error){
        return {error};
    }
}

