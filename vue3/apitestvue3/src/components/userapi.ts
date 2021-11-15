import axios from 'axios';

declare module '@/components/userapi'

const axiosClient = axios.create({
    baseURL: 'http://localhost:1234',
});

export async function getAllUsers(){
    try{
        const {data} = await axiosClient.get(`/users`)
        return [null, data];
    }
    catch(error){
        return [error];
    }
}

export async function getUser(userId: number){
    try{
        const {data} = await axiosClient.get(`users/${userId}`)
        return [data];
    }
    catch(error){
        return [error];
    }
}

export async function addUser(user: JSON){
    try{
        await axiosClient.post(`/usersadd`,user)
        return [null];
    }
    catch(error){
        return [error];
    }
}

export async function editUser(user: JSON, userId: number){
    try{
        await axiosClient.put(`/users/edit/${userId}`,user)
        return [null];
    }
    catch(error){
        return [error];
    }
}

export async function deleteUser(userId: number){
    try{
        await axiosClient.delete(`/users/delete/${userId}`)
        return [null];
    }
    catch(error){
        return [error];
    }
}

