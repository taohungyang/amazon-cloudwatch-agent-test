// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT

//go:build linux && integration
// +build linux,integration

package metric_append_dimension

import (
    "github.com/aws/amazon-cloudwatch-agent-test/environment"
    "github.com/aws/amazon-cloudwatch-agent-test/test/metric"
    "github.com/aws/amazon-cloudwatch-agent-test/test/status"
    "github.com/aws/amazon-cloudwatch-agent-test/test/test_runner"

    "log"
    "time"
)

type AppendDimensionsTestRunner struct {
    test_runner.BaseTestRunner
}

var _ test_runner.ITestRunner = (*AppendDimensionsTestRunner)(nil)

func (t *AppendDimensionsTestRunner) Validate() status.TestGroupResult {
    metricsToFetch := t.GetMeasuredMetrics()
    testResults := make([]status.TestResult, len(metricsToFetch))
    for i, metricName := range metricsToFetch {
        testResults[i] = t.validateNoAppendDimensionMetric(metricName)
    }

    return status.TestGroupResult{
        Name:        t.GetTestName(),
        TestResults: testResults,
    }
}

func (t *AppendDimensionsTestRunner) GetTestName() string {
    return "AppendDimension"
}

func (t *AppendDimensionsTestRunner) GetAgentConfigFileName() string {
    return "append_dimension.json"
}

func (t *AppendDimensionsTestRunner) GetAgentRunDuration() time.Duration {
    return 3 * time.Minute
}

func (t *AppendDimensionsTestRunner) GetMeasuredMetrics() []string {
    return []string{"cpu_time_active"}
}

func (t *AppendDimensionsTestRunner) validateNoAppendDimensionMetric(metricName string) status.TestResult {
    testResult := status.TestResult{
        Name:   metricName,
        Status: status.FAILED,
    }

    fetcher := metric.MetricValueFetcher{
        Env: &environment.MetaData{},
        ExpectedDimensionNames: []string{"ImageId", "InstanceId", "InstanceType"}}

    values, err := fetcher.Fetch("MetricAppendDimensionTest", metricName, metric.AVERAGE)
    log.Printf("metric values are %v", values)
    if err != nil {
        return testResult
    }

    if !isAllValuesGreaterThanOrEqualToZero(metricName, values) {
        return testResult
    }

    // TODO: Range test with >0 and <100
    // TODO: Range test: which metric to get? api reference check. should I get average or test every single datapoint for 10 minutes? (and if 90%> of them are in range, we are good)

    testResult.Status = status.SUCCESSFUL
    return testResult
}

