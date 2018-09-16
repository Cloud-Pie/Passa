package server

const serverTemplate = `

<!DOCTYPE HTML>
<html>
<head>
  <title>Timeline | PASSA</title>

  <style>
    body, html {
      font-family: arial, sans-serif;
      font-size: 11pt;
    }
    input {
      margin: 2px 0;
    }

   table {
    display: table;
    border-collapse: separate;
    border-spacing: 2px;
    border-color: gray;
}
  </style>

<script src="https://ajax.googleapis.com/ajax/libs/jquery/3.3.1/jquery.min.js"></script>

  <script src="https://cdnjs.cloudflare.com/ajax/libs/vis/4.21.0/vis.min.js"></script>
  <link href="https://cdnjs.cloudflare.com/ajax/libs/vis/4.21.0/vis.min.css" rel="stylesheet" type="text/css" />

</head>
<body>

<p>
  <input type="button" id="fit" value="Fit all items"><br>
  <input type="button" id="currentSelection" value="Focus current selection"><br>
</p>


<div id="visualization"></div>

<h4 id="state"></h4>
<table id="myList" border = "1" cellpadding = "5" cellspacing = "5"></table>

</table>
<script>
  // create a dataset with items
  // we specify the type of the fields 'start' and 'end' here to be strings
  // containing an ISO date. The fields will be outputted as ISO dates
  // automatically getting data from the DataSet via items.get().
  var items = new vis.DataSet({
    type: { start: 'ISODate', end: 'ISODate' }
  });



  var container = document.getElementById('visualization');
  var options = {
    //start: '2014-01-10',
    //end: '2014-02-10',
    editable: true,
    showCurrentTime: true
  };


  document.getElementById('fit').onclick = function() {
    timeline.fit();
  };
  document.getElementById('currentSelection').onclick = function() {
    var selection = timeline.getSelection();
    timeline.focus(selection);
  };



function fixTimeline(timeline){
    timeline.on('select', function (properties) {

      selectedState=items.get(properties.items[0]);
      console.log(selectedState);
      title=document.getElementById("state");
      title.innerText=selectedState.content;

      myList=document.getElementById("myList");
      myList.innerHTML=" \
    <tr> \
    <th>Name</th> \
    <th>Scale</th> \
    </tr> \
    "

    Object.keys(selectedState.services)
      .forEach(function eachKey(key) {
        console.log(key); // alerts key
        console.log(selectedState.services[key]); // alerts value

        var row=myList.insertRow(1);
        var cell1 = row.insertCell(0);
        var cell2 = row.insertCell(1);

        cell1.innerHTML = key;
        cell2.innerHTML = selectedState.services[key].Replicas;
      });


    });
}


function getStates(){
  console.log("get states called")
  $.ajax(
      {
          type: "get",
          url: "/api/states",
          success: function (data) {
              data.forEach(element => {
                new_item={
                    start:element.ISODate,
                    content: element.Name,
                    services: element.Services
                }
                console.log(new_item)
                items.add(new_item)


              });

  var timeline = new vis.Timeline(container, items, options); //Set the timeline
  fixTimeline(timeline)
          }
      }
  )
}
getStates()
</script>

</body>
</html>

`
