<template>
  <div>

    <div class="box box-default">
      <div class="box-header">
        <h3 class="box-title">Diagram: <b>{{ obj.namespace }} / {{ obj.kind }} / {{ obj.name }}</b></h3>
      </div>

      <div class="overlay" v-if="loading">
        <i class="fa fa-refresh fa-spin"></i>
      </div>
      <div class="box-body">
        <div >
          <div class="col-xs-12">
            <div id="object_details_graph" v-bind:style="{ height: height }"></div>
          </div>
        </div>
        <!-- /.row -->
      </div>
    </div>

  </div>
</template>
<script>
  import { getObjectDiagram } from 'lib/api.js'
  import vis from 'vis'

  export default {
    data () {
      // empty data
      return {
        loading: false,
        error: null
      }
    },
    created () {
      // fetch the data when the view is created and the data is already being observed
      this.fetchData()
    },
    props: ['obj', 'height'],
    methods: {
      fetchData () {
        this.loading = true
        this.error = null

        const fetchSuccess = $.proxy(function (data) {
          this.loading = false
          let container = document.getElementById('object_details_graph')
          let network = new vis.Network(container, data, options)
          network.fit()
        }, this)

        const fetchError = $.proxy(function (err) {
          this.loading = false
          this.error = err
        }, this)

        getObjectDiagram(this.obj, fetchSuccess, fetchError)
      }
    }
  }

  var options = {
    nodes: {
      font: {
        size: 16,
        color: 'white'
      },
      borderWidth: 0,
      chosen: {
        label: false,
        node: chosenNode
      }
    },
    edges: {
      width: 1,
      font: {
        size: 16,
        strokeWidth: 0,
        color: 'rgb(246,65,111)',
        align: 'top'
      }
    },
    groups: {
      bundle: {
        size: 25,
        shape: 'circularImage',
        image: '/static/img/bundle-icon-circle.png',
        color: {background: 'white', border: 'lightgray'}
      },
      service: {
        font: {
          color: 'rgb(220,213,31)',
          multi: 'html'
        },
        size: 25,
        shape: 'circularImage',
        image: '/static/img/service-icon.png',
        color: {background: 'rgb(164,253,74)', border: 'rgb(220,213,31)'}
      },
      componenthelm: {
        font: {
          color: 'rgb(66,136,251)',
          multi: 'html'
        },
        size: 25,
        shape: 'circularImage',
        image: '/static/img/helm-icon.png',
        color: {background: 'white', border: 'rgb(66,136,251)'}
      },
      componentraw: {
        font: {
          color: 'rgb(66,136,251)',
          multi: 'html'
        },
        size: 25,
        shape: 'circularImage',
        image: '/static/img/k8s-icon.png',
        color: {background: 'white', border: 'rgb(66,136,251)'}
      },
      bundleInstance: {
        font: {
          color: 'black',
          multi: 'html'
        },
        color: {background: 'rgb(19,132,186)', border: 'darkslategrey'},
        shape: 'box'
      },
      claim: {
        size: 25,
        shape: 'circularImage',
        image: '/static/img/user-icon-circle.png',
        color: {background: 'white', border: 'lightgray'}
      },
      claimNotResolved: {
        size: 25,
        shape: 'circularImage',
        image: '/static/img/user-icon-circle.png',
        color: {background: 'red', border: 'white'}
      },
      error: {
        shape: 'icon',
        font: {
          multi: 'html'
        },
        icon: {
          face: 'FontAwesome',
          code: '\uf235',
          size: 50,
          color: 'red'
        }
      }
    },
    layout: {
      randomSeed: 239,
      hierarchical: {
        direction: 'LR',
        levelSeparation: 220
      }
    },
    interaction: {
      hover: true,
      navigationButtons: true,
      keyboard: true
    },
    physics: false
  }

  function chosenNode (values, id, selected, hovering) {
    values.color = '#ffdd88'
    values.borderColor = '#ff0000'
  }
/*
  function clickedNode (params) {
    params.event = "[original event]";
    var node = this.getNodeAt(params.pointer.DOM);
    var edge = this.getEdgeAt(params.pointer.DOM);

    var id = "";
    if (node) {
      id = node
    } else if (edge) {
      id = edge
    }
    if (id) {
      // loadJSON("/api/object-view?id=" + id, objectLoaded, objectNotLoaded);
      // $("#rule-log-button").click();
    } else {
      // app.obj_view = [];
    }
  }
*/
</script>
<style>
  #object_details_graph {
    color: #d3d3d3;
    background-color: #222222;
    border: 1px solid #444444;
    font: 16pt arial;
    width: 100%;
  }
</style>
