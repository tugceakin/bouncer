/**
 * Created by tugceakin on 4/8/16.
 */

bouncerApp.controller('BenchmarkingController', function($scope, $interval, $http, $parse, benchmarking, configuration) {
    $scope.pageClass = 'page-benchmarking';

    $scope.graphOff = true;
    $scope.benchmarkCompleted = false;
    $scope.benchmarkInput = "";
    $scope.stats = {};
    $scope.statsShown = false;
    $scope.currentConfigId = 1;
    $scope.backendServers = [];
    $scope.backendServerCount = 1;

    benchmarking.setSocketConnection;
    benchmarking.resetGraph($scope);
    benchmarking.onGraphLineClick($scope);


    $scope.startBenchmarking = function(){
        var config = {
            "id": $scope.configId,
            "host": $scope.hostName,
            "path": $scope.path,
            "reqPerSecond": $scope.reqPerSecond,
            "concurrency": $scope.concurrency,
            "backendServers": configuration.getBackendServers($scope)
        };

        console.log(config);
        benchmarking.resetGraph($scope);
        benchmarking.updateGraph($scope);
        //$scope.stats = benchmarking.getBenchmarkingStats();
    }

    $scope.closeConnection = function(){
        benchmarking.closeConnection();
        $scope.graphOff = true;
    }

    $scope.getBenchmarkingStats = function(){
        if($scope.statsShown == false){
            $scope.statsShown = true;
        }else{
            $scope.statsShown = false;
        }
    }


});