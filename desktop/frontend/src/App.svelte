
<main>
  <div class="result" id="result">{resultText}</div>
  <div class="input-box" id="input">
    <input autocomplete="off" bind:value={name} class="input" id="name" type="text"/>
    <!--<button class="btn" on:click={drawCandles}>Greet</button>-->
  </div>
  <div id="container"></div>
</main>

<style>

  #logo {
    display: block;
    width: 50%;
    height: 50%;
    margin: auto;
    padding: 10% 0 0;
    background-position: center;
    background-repeat: no-repeat;
    background-size: 100% 100%;
    background-origin: content-box;
  }

  .result {
    height: 20px;
    line-height: 20px;
    margin: 1.5rem auto;
  }

  .input-box .btn {
    width: 60px;
    height: 30px;
    line-height: 30px;
    border-radius: 3px;
    border: none;
    margin: 0 0 0 20px;
    padding: 0 8px;
    cursor: pointer;
  }

  .input-box .btn:hover {
    background-image: linear-gradient(to top, #cfd9df 0%, #e2ebf0 100%);
    color: #333333;
  }

  .input-box .input {
    border: none;
    border-radius: 3px;
    outline: none;
    height: 30px;
    line-height: 30px;
    padding: 0 10px;
    background-color: rgba(240, 240, 240, 1);
    -webkit-font-smoothing: antialiased;
  }

  .input-box .input:hover {
    border: none;
    background-color: rgba(255, 255, 255, 1);
  }

  .input-box .input:focus {
    border: none;
    background-color: rgba(255, 255, 255, 1);
  }

</style>

<script lang="ts">
  import logo from './assets/images/logo-universal.png'
  import {InitCandles} from '../wailsjs/go/main/App.js'
  import * as echarts from 'echarts'
  import {EventsOn} from "../wailsjs/runtime";
  import Highcharts from 'highcharts/highstock';
  import dayjs from 'dayjs';
  let resultText: string = "Please enter your name below 👇. Yes Mother Fucker"
  let candles = []
  let name: string

  let chart;
  let isReady = false
  EventsOn("candles.init",result => {
    candles = result
    if (isReady == false) {
      drawChart()
    }
    isReady = true
  })
  EventsOn("candles.new",result => {
    let res = JSON.parse(result)
    candles = [...candles, res]
    if (res != undefined || res != null) {
      drawChart()
    }
  })

  function drawChart() {
    try {

      var option;
      let chartData = []
      const upColor = '#ec0000';
      const upBorderColor = '#8A0000';
      const downColor = '#00da3c';
      const downBorderColor = '#008F28';
      for (let i = 0; i<candles.length; i++) {
        let candle = candles[i]
        let time = dayjs(new Date(candle["t"]).getTime()).format("dddd, MMM D, hh:mm a")
        let data = [
          time,
          candle["o"],
          candle["h"],
          candle["l"],
          candle["c"],
        ]
        chartData.push(data)
      }

      // create the chart
      chart = Highcharts.chart('container', {
        rangeSelector: {
          selected: 1
        },
        xAxis:{
          zoomEnabled:true,
        },
        yAxis:{
          zoomEnabled:true,
        },
        accessibility:{
          enabled:true
        },
        title: {
          text: "XRPUSDT"
        },
        annotations: {
          draggable: "xy"
        },
        chart:{
          zooming:{
            key:"ctrl",
            pinchType:"x",
            singleTouch:true,
            mouseWheel:true,
            type:"x"
          }
        },
        navigator:{
          enabled:true
        },
        series: [{
          type: 'candlestick',
          name: 'XRPUSDT',
          data: chartData,
          dataGrouping: {
            units: [
              [
                'week', // unit name
                [1] // allowed multiples
              ], [
                'month',
                [1, 2, 3, 4, 6]
              ]
            ]
          }
        }]
      });

    }catch (err) {
      alert(err)
    }

  }
</script>
