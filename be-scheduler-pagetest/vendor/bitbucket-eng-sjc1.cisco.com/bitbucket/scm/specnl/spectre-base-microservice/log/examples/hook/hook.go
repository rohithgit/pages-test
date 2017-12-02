package hook // spectre

// NOTE: This example has build errors because airbrake uses logrus
// and this is a fork of logrus. It has therefore been commented out
// to allow for a clean build. The code can still be used as an
// example of how to use hooks.

// package main

/*
 *import (
 *    "bitbucket-eng-sjc1.cisco.com/bitbucket/specnl/spectre-base-microservice/log"
 *    airbrake "gopkg.in/gemnasium/logrus-airbrake-hook.v2"
 *)
 *
 *var logger = log.New()
 *
 *func init() {
 *    logger.Formatter = new(log.TextFormatter) // default
 *    logger.Hooks.Add(airbrake.NewHook(123, "xyz", "development"))
 *}
 *
 *func main() {
 *    logger.WithFields(log.Fields{
 *        "animal": "walrus",
 *        "size":   10,
 *    }).Info("A group of walrus emerges from the ocean")
 *
 *    logger.WithFields(log.Fields{
 *        "omg":    true,
 *        "number": 122,
 *    }).Warn("The group's number increased tremendously!")
 *
 *    logger.WithFields(log.Fields{
 *        "omg":    true,
 *        "number": 100,
 *    }).Fatal("The ice breaks!")
 *}
 */
