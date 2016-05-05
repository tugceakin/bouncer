/**
 * Created by tugceakin on 4/8/16.
 */

bouncerApp.controller('BenchmarkingController', function($scope, $interval, $http, graphFactory) {
    $scope.pageClass = 'page-benchmarking';

    $scope.graphOff = true;
    $scope.benchmarkCompleted = false;
    $scope.benchmarkInput = "";
    $scope.stats = {};
    $scope.statsShown = false;

    graphFactory.resetGraph($scope);
    graphFactory.onGraphLineClick($scope);


    $scope.startBenchmarking = function(){
        graphFactory.resetGraph($scope);
        graphFactory.setBenchmarkingStats($scope);
        //$scope.stats = graphFactory.getBenchmarkingStats();
    }

    $scope.getBenchmarkingStats = function(){
        if($scope.statsShown == false){
            $scope.statsShown = true;
        }else{
            console.log($scope.statsShown)
            $scope.statsShown = false;
        }
        console.log($scope.stats.server_port);
    }

});


// bouncerApp.controller("StatsController", function ($scope, graphFactory) {
//   console.log("bar controller");

//   graphFactory.setBenchmarkingStats($scope);
//   console.log($scope.stats);
//   $scope.labels = ['Min', 'Mean', '[+/-sd]', 'Median', 'Max'];
//   $scope.series = ['Connect', 'Processing', 'Waiting'];

//   $scope.data = [
//     [2,3,4,5],
//     [4,5,6,7],
//     [8,9,10,11],
//     [8,9,10,11]
//   ];
// });
       

