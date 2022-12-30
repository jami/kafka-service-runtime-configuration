<script setup>
import EditActionIcon from './icons/schemaEditAction.vue'
</script>

<template>
  <div class="schema_list_item">
    <div class="schema_header">
      <h2>Application ID - {{ schema.id }}</h2>
    </div>
    <div class="schema_action ">
      <a href="#" v-if="!showEditor" class="toggleButton" @click="showEditor = !showEditor">+</a>
      <a href="#" v-if="showEditor" class="toggleButton" @click="showEditor = !showEditor">-</a>
    </div>
    <div v-if="showEditor" class="schema_view">
      <vue-json-pretty :data="schema.schema" />
    </div>
    <div v-if="showEditor" class="schema_value">
      <vue-json-editor 
        v-model="json" 
        mode="code"
        :modes="['code']"
        :show-btns="false"
        :exapndedOnStart="true"
        @json-change="onJsonChange"></vue-json-editor>
    </div>
    <div v-if="showEditor" class="schema_footer">
      <a href="#" v-if="contentChanged" class="saveValuesButton" @click="saveValues">Update</a>
    </div>
  </div>
</template>

<style scoped>
.schema_list_item {
  display: grid;
  grid-template-areas: 
    'header header header header header header_action'
    'left left left right right right'
    'footer footer footer footer footer footer';
  padding-bottom: 3vh;
}

.schema_header {
  grid-area: header;
  padding: 10px;
}

.schema_action {
  grid-area: header_action;
  display: flex;
  justify-content: flex-end;
  margin-left: auto;
  margin-right: 0;
}


.schema_view {
  grid-area: left;
  padding-top: 3em;
}

.schema_value {
  grid-area: right;
  padding-top: 3em;
}

.schema_footer {
  grid-area: footer;
}

.saveValuesButton {
  background-color:#44c767;
	border-radius:4px;
	border:1px solid #18ab29;
	display:inline-block;
	cursor:pointer;
	color:#ffffff;
	font-family:Arial;
	font-size:17px;
	padding:6px 6px;
	text-decoration:none;
	text-shadow:0px 1px 0px #2f6627;
}

.toggleButton {
	background:linear-gradient(to bottom, #33bdef 5%, #019ad2 100%);
	background-color:#33bdef;
	border-radius:4px;
	border:1px solid #057fd0;
	display:inline-block;
	cursor:pointer;
	color: #ffffff;
	font-family:Arial;
	font-size:19px;
	font-weight:bold;
	padding:6px 6px;
	text-decoration:none;
  aspect-ratio: 1 / 1;
  width: auto;
  height: 80%;
  text-align: center;
}

.toggleButton:hover {
	background:linear-gradient(to bottom, #019ad2 5%, #33bdef 100%);
	background-color:#019ad2;
}

.toggleButton:active {
	position:relative;
	top:1px;
}
</style>

<script>
import VueJsonEditor from 'vue-json-editor'
import VueJsonPretty from "vue-json-pretty";
import "vue-json-pretty/lib/styles.css";

export default {
  data () {
    console.log("json ", this.schema.values)
    return {
      contentChanged: false,
      showEditor: false,
      appId: this.schema.id, 
      json: this.schema.values
    }
  },
  props: {
    schema: {
      type: Object
    }
  },
  components: {
    VueJsonEditor,
    VueJsonPretty,
  },
  watch: {
    schema: () => {
      console.log("watch schema prop change")
      json = this.schema.values
    }
  },
  methods: {
    onJsonChange(value) {
      this.contentChanged = true
      console.log('value:', value)
    },
    saveValues(value) {
      console.log('save values')
      // POST request using fetch with error handling
      const requestOptions = {
        method: 'POST',
        headers: { 
          'Content-Type': 'application/json' 
        },
        body: JSON.stringify({
          id: this.appId,
          values: this.json
        })
      };
      fetch('/api/schema/update', requestOptions)
    .then(async response => {
      const data = await response.json();

      // check for error response
      if (!response.ok) {
        // get error message from body or default to response status
        const error = (data && data.message) || response.status;
        return Promise.reject(error);
      }

      console.log('savevalues res', data, response);
    })
    .catch(error => {
      this.errorMessage = error;
      console.error('There was an error!', error);
    });
    }
  }
};
</script>
