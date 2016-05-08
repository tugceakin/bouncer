bouncerApp.factory('graphFactory', function ($http, $interval) {

    
   return {
       drawCircles: function($scope) {
          var w = 700;
          var h = 3600;
          var barPadding = 2;
          //Create SVG element
          var svg = d3.select(".box")
              .append("svg")
              .attr("width", w)
              .attr("height", h);
         }

   }
});