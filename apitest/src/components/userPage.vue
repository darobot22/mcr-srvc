<template>
    <div id="inputId">
            <input v-model.number="userId" placeholder="user ID">
            <button @click="HandleGetUser">Get user</button> <button @click="EditUserButton">Edit</button> <button @click="HandleDeleteUser">Delete</button> <button @click="AddUserButton">Add</button>
            <p v-if="userGet">ID: {{user.Id}} | Name: {{user.Name}} | Email: {{user.Email}} | Phone: {{user.Phone}}</p>
            <div v-if="editing">
                Name: <input v-model="user.Name" placeholder="Name">
                Email: <input v-model="user.Email" placeholder="Email">
                Phone: <input v-model="user.Phone" placeholder="Phone">
                Password: <input v-model="user.Password" placeholder="Password"><br>
                <button @click="HandleEditUser">Save</button>
            </div>
            <div v-if="add">
                Name: <input v-model="user.Name" placeholder="Name">
                Email: <input v-model="user.Email" placeholder="Email">
                Phone: <input v-model="user.Phone" placeholder="Phone">
                Password: <input v-model="user.Password" placeholder="Password"><br>
                <button @click="HandleAddUser">Save</button>
            </div>
                <div id="cat">
           <div v-for="serv in services" v-bind:key="serv.Id">
                <h2>{{serv.Name}}</h2>
                <p>{{serv.Code}}</p>
                <button @click="Request(serv.Id)">Request</button>
           </div>
    </div>
    </div>
</template>

<script>
import axios from 'axios';
const axiosClient = axios.create({
    baseURL: 'http://localhost:1234',
});
export default {
    name: 'userPage',
    data(){
        return{
            user: JSON,
            userId: null,
            editing: false,
            userGet: false,
            add: false,
            services: [],
        };
    },

    methods: {
        HandleGetUser(){
            this.userGet = true
            axiosClient.get(`/users/${this.userId}`)
            .then(response => {
                this.user = response.data
                        this.HandleGetServices()

            })
            .catch(e => {
                this.errors.push(e)
            })
        },
        HandleEditUser(){
            this.editing = false
            axiosClient.put(`/users/edit/${this.userId}`, this.user)
            .catch(e => {
                this.errors.push(e)
            })
        },
        EditUserButton(){
            this.editing = true
        },
        HandleDeleteUser(){
            axios.delete(`/users/delete/${this.userId}`)
        },
        AddUserButton(){
            this.add = true
        },
        HandleAddUser(){
            this.add = false
            this.user.Id = this.userId
            axios.post(`/usersadd`,this.user)
        },
        HandleGetServices(){
            this.userGet = true
            this.services = []
            axios.get(`http://localhost:1235/services/`)
            .then(response => {
                console.log(response.data.split('\n'))
                    response.data.split('\n').map(s=>this.services.push(JSON.parse(s)))
                })
            .catch(e => {
                console.error(e)
            })
        },
        HandleEditService(){
            this.sediting = false
            axios.put(`http://localhost:1235/services/edit/${this.serviceId}`, this.service)
            .catch(e => {
                this.errors.push(e)
            })
        },
        EditServiceButton(){
            this.sediting = true
        },
        HandleDeleteService(){
            axios.delete(`http://localhost:1235/services/delete/${this.serviceId}`)
        },
        AddServiceButton(){
            this.sadd = true
        },
        HandleAddService(){
            this.sadd = false
            this.service.Id = this.serviceId
            axios.post(`http://localhost:1235/servicesadd`,this.service)
        },
        Request(sid){
            console.log(sid)
            axios.post(`http://localhost:1235/services/request/${sid}?uid=${this.userId}`)
        }
    },
}
</script>