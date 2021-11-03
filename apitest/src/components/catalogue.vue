<template>
    <div id="serv">
           <div v-for="serv in services" v-bind:key="serv.id">
                <h2>{{serv.Name}}</h2>
                <p>{{serv.Code}}</p>
                <button @click="Request(serv.id)">Request</button>
           </div>
    </div>
</template>

<script>
import {getAllServices,addService,editService,deleteService} from '@/components/catalogueapi.js'
import axios from 'axios'
import {userId} from '@/components/userPage3.vue'

export default {
    name: 'catalogue',
    data(){
        return{
            sediting: false,
            sadd: false,
            services: [],
        };
    },

    methods: {
        HandleGetServices(){
            const [res, error] = getAllServices()
            if (error) console.error(error)
            else this.services = res    
        },
        HandleEditService(service){
            this.sediting = false
            const [error] = editService(service)
            if (error) console.error(error)
        },
        EditServiceButton(){
            this.sediting = true
        },
        HandleDeleteService(serviceId){
            const [error] = deleteService(serviceId)
            if (error) console.error(error)
        },
        AddServiceButton(){
            this.sadd = true
        },
        HandleAddService(service){
            this.add = false
            const [error] = addService(service)
            if (error) console.error(error)
        },
        Request(sid){
            axios.post(`http://localhost:1235/services/request/${sid}?uid=${userId}`)
        }
    },
    
    created(){
        this.HandleGetServices()
    }
}
</script>