import {getUser, addUser, deleteUser, editUser} from "./userapi.ts";

var inputId = new Vue({
        el: '#inputId',
        data: {
            userId : null,
            user: null
        },
        methods: {
            HandleGetUser(){
                const [error, user] = await getUser(this.userId)
                if(error) console.error(error);
                else this.user = user;
            },
            HandleUpdateClick(){
                const [error, editedUser] = await editUser(this.userId)
                if (error) console.error(error);
                else this.user = editedUser;
            },
            HandleAddClick(){
                const [error, addedUser] = await addUser()
                if (error) console.error(error);
                else this.user = addedUser;
            },
            HandleDelete(){
                const [error] = await deleteUser(this.userId)
                if (error) console.error(error);
            },
        }
})