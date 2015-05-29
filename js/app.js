var mseApp = angular.module('mseApp', ['ngMaterial']);

mseApp.controller('mseCtrl', function($scope, $http, $mdSidenav){
    
    $scope.status = [];
    
    $scope.getBoard = function() {
        $http.get('/api/board', $scope.cfg).success(function(d){
            $scope.board = d;
            if (d.State == "End") {
                return;
            }
    
            $scope.metalTrack = {
                cls: "metal",
                name: "Metal Storage", 
                value: d.MetalStorage};
            $scope.wealthTrack = {
                cls: "wealth",
                name: "Wealth Storage",
                value: d.WealthStorage
            };
            $scope.militaryTrack = {
                cls: "military",
                name: "Military Strength",
                value: d.MilitaryStrength
            };
            
            if ($scope.board.state == "End") {
                return;
            }
            return $scope.getBoard();
        });
    };
    
    $scope.getStatus = function() {
        $http.get('/api/status', $scope.cfg).success(function(d) {
            $scope.status.push(d);
            if (d.End) {
                return
            }
            return $scope.getStatus();
        })
    };
    
    $scope.getPrompt = function() {
        $http.get('/api/prompt', $scope.cfg).success(function(d) {
            $scope.prompt = d;
            if (d.End) {
                return
            }
            return $scope.getPrompt();
        })
    };

    $scope.makeChoice = function(key) {
        $http.post('/api/choice', {ID: $scope.cfg.params.ID, Key: key})
            .success(function(d){
            })
            .error(function(d){
                $scope.status.push('Choice failed: ' + d);
            });
    };

    $scope.showBoard = function() {
        $mdSidenav('board').toggle();
    };
    
    $scope.hideBoard = function() {
        $mdSidenav('board').close();
    };

    $scope.showStatus = function() {
        $mdSidenav('status').toggle();
    };
    
    $scope.hideStatus = function() {
        $mdSidenav('status').close();
    };

    $http.get('/api/newGame').success(function(d){
        $scope.cfg = {params: {ID: d.ID}};
        $scope.getBoard();
        $scope.getStatus();
        $scope.getPrompt();
    });
    
});