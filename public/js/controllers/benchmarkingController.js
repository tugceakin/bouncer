/**
 * Created by tugceakin on 4/8/16.
 */

bouncerApp.controller('BenchmarkingController', function($scope, $interval, $http, $parse, benchmarking, configuration) {
    $scope.graphOff = true;
    $scope.benchmarkCompleted = false;
    $scope.benchmarkInput = "";
    $scope.stats = {};
    $scope.statsShown = false;
    $scope.currentConfigId = 1;
    $scope.backendServers = [];

    benchmarking.setSocketConnection;
    benchmarking.resetGraph($scope);
    benchmarking.onGraphLineClick($scope);

    configuration.getAllConfigs($scope);

    $scope.startBenchmarking = function(){
        benchmarking.resetGraph($scope);
        benchmarking.updateGraph($scope);
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