bouncerApp.factory('graphFactory', function ($http, $interval) {

    
    var benchmarkingStats = {};

    var stopBenchmarking = function($scope){
            if (angular.isDefined($scope.Benchmarking)) {
                $interval.cancel($scope.Benchmarking);
                $scope.Benchmarking = undefined;
                timeCounter = 0;
                $scope.graphOff = true;
                $scope.benchmarkCompleted = true;
                //Hardcoded for now
                benchmarkingStats = {
                  "server_hostname" :      "localhost",
                  "server_port":           "9090",
                  "document_path":         "/sdf",
                  "concurrency_level":     "100",
                  "time_take_for_tests":   "6.102 seconds",
                  "complete_requests":      "10000",
                  "failed_requests":        "0",
                  "total_transferred":      "1400000 bytes",
                  "html_transferred":       "230000 bytes",
                  "requests_per_sec":    "1638.90 [#/sec] (mean)",
                  "time_per_request":       "61.017 [ms] (mean)",
                  "time_per_request":       "0.610 [ms] (mean, across all concurrent requests)",
                  "transfer_rate":          "224.07 [Kbytes/sec] received",
                  "connection_times": {
                    "connect": {
                      "min": 0,
                      "mean": 3,
                      "[+/-sd]": 4.3,
                      "median": 1,
                      "max": 60
                    },
                    "processing": {
                      "min": 50,
                      "mean": 58,
                      "[+/-sd]": 10.0,
                      "median": 54,
                      "max": 229
                    },
                    "waiting": {
                      "min": 50,
                      "mean": 57,
                      "[+/-sd]": 9.9,
                      "median": 54,
                      "max": 229
                    }
                  }
                }
            }
        }



   return {
       onGraphLineClick: function($scope) {
          $scope.onClick = function (points, evt) {
              console.log(points, evt);
          };
       },

       resetGraph: function($scope) {
          $scope.labels =['', '', '', '', '', ''];
          $scope.data = [
              [0,0,0,0,0,0]
          ];
       },

       setBenchmarkingStats: function($scope){
          var timeCounter = 0;

          //Use this request later to connect graphs to the actual results.
          $http({
              method: 'POST',
              url: "/startBenchmarking",
              headers: {'Content-Type': 'application/x-www-form-urlencoded'}, 
              data: {benchmarkInput: $scope.benchmarkInput}
          }).success(function (data) {
              console.log(data);
          });

          $scope.benchmarkCompleted = false;
          $scope.graphOff = false;
          // Simulate async data update
          this.resetGraph($scope);
          $scope.Benchmarking = $interval(function () {
              timeCounter += 3;
              // Remove first element
              $scope.labels.splice(0,1);
              $scope.data[0].splice(0,1);

              $scope.data[0].push(Math.random() * 1000); //For now get random data and push it
              $scope.labels.push(timeCounter);


              if(timeCounter >= 21) {
                  stopBenchmarking($scope);
                  $scope.stats = benchmarkingStats;
              }

          }, 30);
       },

       getBenchmarkingStats: function($scope){
        return benchmarkingStats;
       }

   }
});