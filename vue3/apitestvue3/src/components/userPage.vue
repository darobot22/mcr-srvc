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
            <div v-if="userGet">
                <h1>Service history</h1>
                <div v-for="his in history" v-bind:key="his.Id">
                    <h2>{{his.ServiceName}}</h2>
                    <p>{{his.ServiceCode}}</p>
                    <p>{{his.CreateDate}}</p>
                    <p>{{his.ResultData}}</p>
                    <p>{{his.ExecutionDate}}</p>
                </div>
                <h1>Service list</h1>
                <div v-for="serv in services" v-bind:key="serv.Id">
                    <h2>{{serv.Name}}</h2>
                    <p>{{serv.Code}}</p>
                    <button @click="Request(serv.Id)">Request</button>
                </div>
           </div>
    </div>
</template>

<script lang="ts">
import { defineComponent} from "vue"
import {getUser,addUser,editUser,deleteUser} from '../components/userapi';
import {getHistory} from '../components/servicehistoryapi';
import {getAllServices,requestService} from '../components/catalogueapi';

export default defineComponent({
    name: 'userPage',
    data(){
        return{
            user: JSON,
            userId: null,
            editing: false,
            userGet: false,
            add: false,
            history: [],
            services: []
        };
    },

    methods: {
        async HandleGetHistory(Id: number){
            this.history = []
            const [error, res] = await getHistory(Id)
            if (error) console.error(error)
            else {
                res.split('\n').map(s=>this.history.push(JSON.parse(s)))  
            }
        },
        async HandleGetServices(){
            this.services = []
            const [error, res] = await getAllServices()
            if (error) console.error(error)
            else {
                console.log(res)
                res.split('\n').map(s=>this.services.push(JSON.parse(s)))
            }
        },
        async HandleGetUser(){
            this.userGet = true
            const [res, error] = await getUser(this.userId)
            if (error) console.error(error)
            else {
            this.user = res
            this.HandleGetHistory(this.user.Id)
            this.HandleGetServices()
            }
        },
        async HandleEditUser(){
            this.editing = false
            const [error] = await editUser(this.user, this.userId)
            if (error) console.error(error)
        },
        EditUserButton(){
            this.editing = true
        },
        async HandleDeleteUser(){
            const [error] = await deleteUser(this.userId)
            if (error) console.error(error)
        },
        AddUserButton(){
            this.add = true
        },
        async HandleAddUser(){
            this.add = false
            this.user.Id = this.userId
            const [error] = await addUser(this.user)
            if (error) console.error(error)
        },
        async Request(sid: number){
            const [error] = await requestService(sid,this.userId)
            if (error) console.error(error)
        }
    }
})
</script> 