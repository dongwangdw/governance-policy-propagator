// Copyright (c) 2020 Red Hat, Inc.

package e2e

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	appsv1 "github.com/open-cluster-management/governance-policy-propagator/pkg/apis/apps/v1"
	policiesv1 "github.com/open-cluster-management/governance-policy-propagator/pkg/apis/policies/v1"
	"github.com/open-cluster-management/governance-policy-propagator/pkg/controller/common"
	"github.com/open-cluster-management/governance-policy-propagator/test/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const case1PolicyName string = "case1-test-policy"
const case1PolicyYaml string = "../resources/case1_propagation/case1-test-policy.yaml"

var _ = Describe("Test policy propagation", func() {
	Describe("Create policy/pb/plc in ns:"+testNamespace+" and then update pb/plc", func() {
		It("should be created in user ns", func() {
			By("Creating " + case1PolicyYaml)
			utils.Kubectl("apply",
				"-f", case1PolicyYaml,
				"-n", testNamespace)
			plc := utils.GetWithTimeout(clientHubDynamic, gvrPolicy, case1PolicyName, testNamespace, true, defaultTimeoutSeconds)
			Expect(plc).NotTo(BeNil())
		})
		It("should propagate to cluster ns managed1", func() {
			By("Patching test-policy-plr with decision of cluster managed1")
			plr := utils.GetWithTimeout(clientHubDynamic, gvrPlacementRule, case1PolicyName+"-plr", testNamespace, true, defaultTimeoutSeconds)
			plr.Object["status"] = utils.GeneratePlrStatus("managed1")
			plr, err := clientHubDynamic.Resource(gvrPlacementRule).Namespace(testNamespace).UpdateStatus(plr, metav1.UpdateOptions{})
			Expect(err).To(BeNil())
			plc := utils.GetWithTimeout(clientHubDynamic, gvrPolicy, testNamespace+"."+case1PolicyName, "managed1", true, defaultTimeoutSeconds)
			Expect(plc).ToNot(BeNil())
			opt := metav1.ListOptions{LabelSelector: common.RootPolicyLabel + "=" + testNamespace + "." + case1PolicyName}
			utils.ListWithTimeout(clientHubDynamic, gvrPolicy, opt, 1, true, defaultTimeoutSeconds)
		})
		It("should propagate to cluster ns managed2", func() {
			By("Patching test-policy-plr with decision of cluster managed2")
			plr := utils.GetWithTimeout(clientHubDynamic, gvrPlacementRule, case1PolicyName+"-plr", testNamespace, true, defaultTimeoutSeconds)
			plr.Object["status"] = utils.GeneratePlrStatus("managed2")
			plr, err := clientHubDynamic.Resource(gvrPlacementRule).Namespace(testNamespace).UpdateStatus(plr, metav1.UpdateOptions{})
			Expect(err).To(BeNil())
			plc := utils.GetWithTimeout(clientHubDynamic, gvrPolicy, testNamespace+"."+case1PolicyName, "managed2", true, defaultTimeoutSeconds)
			Expect(plc).ToNot(BeNil())
			opt := metav1.ListOptions{LabelSelector: common.RootPolicyLabel + "=" + testNamespace + "." + case1PolicyName}
			utils.ListWithTimeout(clientHubDynamic, gvrPolicy, opt, 1, true, defaultTimeoutSeconds)
		})
		It("should propagate to cluster ns managed1 and managed2", func() {
			By("Patching test-policy-plr with decision of both managed1 and managed2")
			plr := utils.GetWithTimeout(clientHubDynamic, gvrPlacementRule, case1PolicyName+"-plr", testNamespace, true, defaultTimeoutSeconds)
			plr.Object["status"] = utils.GeneratePlrStatus("managed1", "managed2")
			plr, err := clientHubDynamic.Resource(gvrPlacementRule).Namespace(testNamespace).UpdateStatus(plr, metav1.UpdateOptions{})
			Expect(err).To(BeNil())
			opt := metav1.ListOptions{LabelSelector: common.RootPolicyLabel + "=" + testNamespace + "." + case1PolicyName}
			utils.ListWithTimeout(clientHubDynamic, gvrPolicy, opt, 2, true, defaultTimeoutSeconds)
		})
		It("should remove policy from ns managed1 and managed2", func() {
			By("Patching test-policy-pb with a plc with wrong name")
			pb := utils.GetWithTimeout(clientHubDynamic, gvrPlacementBinding, case1PolicyName+"-pb", testNamespace, true, defaultTimeoutSeconds)
			pb.Object["placementRef"] = &policiesv1.Subject{
				APIGroup: "apps.open-cluster-management.io",
				Kind:     "PlacementRule",
				Name:     case1PolicyName + "-plc-nonexists",
			}
			pb, err := clientHubDynamic.Resource(gvrPlacementBinding).Namespace(testNamespace).Update(pb, metav1.UpdateOptions{})
			Expect(err).To(BeNil())
			opt := metav1.ListOptions{LabelSelector: common.RootPolicyLabel + "=" + testNamespace + "." + case1PolicyName}
			utils.ListWithTimeout(clientHubDynamic, gvrPolicy, opt, 0, true, defaultTimeoutSeconds)
		})
		It("should propagate to cluster ns managed1 and managed2", func() {
			By("Patching test-policy-pb with correct plc")
			pb := utils.GetWithTimeout(clientHubDynamic, gvrPlacementBinding, case1PolicyName+"-pb", testNamespace, true, defaultTimeoutSeconds)
			pb.Object["placementRef"] = &policiesv1.Subject{
				APIGroup: "apps.open-cluster-management.io",
				Kind:     "PlacementRule",
				Name:     case1PolicyName + "-plc",
			}
			pb, err := clientHubDynamic.Resource(gvrPlacementBinding).Namespace(testNamespace).Update(pb, metav1.UpdateOptions{})
			Expect(err).To(BeNil())
			opt := metav1.ListOptions{LabelSelector: common.RootPolicyLabel + "=" + testNamespace + "." + case1PolicyName}
			utils.ListWithTimeout(clientHubDynamic, gvrPolicy, opt, 0, true, defaultTimeoutSeconds)
		})
		It("should remove policy from ns managed1 and managed2", func() {
			By("Patching test-policy-pb with a plc with wrong apigroup")
			pb := utils.GetWithTimeout(clientHubDynamic, gvrPlacementBinding, case1PolicyName+"-pb", testNamespace, true, defaultTimeoutSeconds)
			pb.Object["placementRef"] = &policiesv1.Subject{
				APIGroup: "apps.open-cluster-management.io1",
				Kind:     "PlacementRule",
				Name:     case1PolicyName + "-plc-nonexists",
			}
			pb, err := clientHubDynamic.Resource(gvrPlacementBinding).Namespace(testNamespace).Update(pb, metav1.UpdateOptions{})
			Expect(err).To(BeNil())
			opt := metav1.ListOptions{LabelSelector: common.RootPolicyLabel + "=" + testNamespace + "." + case1PolicyName}
			utils.ListWithTimeout(clientHubDynamic, gvrPolicy, opt, 0, true, defaultTimeoutSeconds)
		})
		It("should propagate to cluster ns managed1 and managed2", func() {
			By("Patching test-policy-pb with correct plc")
			pb := utils.GetWithTimeout(clientHubDynamic, gvrPlacementBinding, case1PolicyName+"-pb", testNamespace, true, defaultTimeoutSeconds)
			pb.Object["placementRef"] = &policiesv1.Subject{
				APIGroup: "apps.open-cluster-management.io",
				Kind:     "PlacementRule",
				Name:     case1PolicyName + "-plc",
			}
			pb, err := clientHubDynamic.Resource(gvrPlacementBinding).Namespace(testNamespace).Update(pb, metav1.UpdateOptions{})
			Expect(err).To(BeNil())
			opt := metav1.ListOptions{LabelSelector: common.RootPolicyLabel + "=" + testNamespace + "." + case1PolicyName}
			utils.ListWithTimeout(clientHubDynamic, gvrPolicy, opt, 0, true, defaultTimeoutSeconds)
		})
		It("should remove policy from ns managed1 and managed2", func() {
			By("Patching test-policy-pb with a plc with wrong kind")
			pb := utils.GetWithTimeout(clientHubDynamic, gvrPlacementBinding, case1PolicyName+"-pb", testNamespace, true, defaultTimeoutSeconds)
			pb.Object["placementRef"] = &policiesv1.Subject{
				APIGroup: "apps.open-cluster-management.io",
				Kind:     "PlacementRule1",
				Name:     case1PolicyName + "-plc-nonexists",
			}
			pb, err := clientHubDynamic.Resource(gvrPlacementBinding).Namespace(testNamespace).Update(pb, metav1.UpdateOptions{})
			Expect(err).To(BeNil())
			opt := metav1.ListOptions{LabelSelector: common.RootPolicyLabel + "=" + testNamespace + "." + case1PolicyName}
			utils.ListWithTimeout(clientHubDynamic, gvrPolicy, opt, 0, true, defaultTimeoutSeconds)
		})
		It("should propagate to cluster ns managed1 and managed2", func() {
			By("Patching test-policy-pb with correct plc")
			pb := utils.GetWithTimeout(clientHubDynamic, gvrPlacementBinding, case1PolicyName+"-pb", testNamespace, true, defaultTimeoutSeconds)
			pb.Object["placementRef"] = &policiesv1.Subject{
				APIGroup: "apps.open-cluster-management.io",
				Kind:     "PlacementRule",
				Name:     case1PolicyName + "-plc",
			}
			pb, err := clientHubDynamic.Resource(gvrPlacementBinding).Namespace(testNamespace).Update(pb, metav1.UpdateOptions{})
			Expect(err).To(BeNil())
			opt := metav1.ListOptions{LabelSelector: common.RootPolicyLabel + "=" + testNamespace + "." + case1PolicyName}
			utils.ListWithTimeout(clientHubDynamic, gvrPolicy, opt, 0, true, defaultTimeoutSeconds)
		})
		It("should remove policy from ns managed1 and managed2", func() {
			By("Patching test-policy-plr with no decision")
			plr := utils.GetWithTimeout(clientHubDynamic, gvrPlacementRule, case1PolicyName+"-plr", testNamespace, true, defaultTimeoutSeconds)
			plr.Object["status"] = &appsv1.PlacementRuleStatus{}
			plr, err := clientHubDynamic.Resource(gvrPlacementRule).Namespace(testNamespace).UpdateStatus(plr, metav1.UpdateOptions{})
			Expect(err).To(BeNil())
			opt := metav1.ListOptions{LabelSelector: common.RootPolicyLabel + "=" + testNamespace + "." + case1PolicyName}
			utils.ListWithTimeout(clientHubDynamic, gvrPolicy, opt, 0, true, defaultTimeoutSeconds)
		})
		It("should clean up", func() {
			utils.Kubectl("delete",
				"-f", case1PolicyYaml,
				"-n", testNamespace)
			opt := metav1.ListOptions{}
			utils.ListWithTimeout(clientHubDynamic, gvrPolicy, opt, 0, false, 1)
		})
	})

	Describe("Create policy/pb/plc in ns:"+testNamespace+" and then update policy", func() {
		It("should be created in user ns", func() {
			By("Creating " + case1PolicyYaml)
			utils.Kubectl("apply",
				"-f", case1PolicyYaml,
				"-n", testNamespace)
			plc := utils.GetWithTimeout(clientHubDynamic, gvrPolicy, case1PolicyName, testNamespace, true, defaultTimeoutSeconds)
			Expect(plc).NotTo(BeNil())
		})
		It("should propagate to cluster ns managed1", func() {
			By("Patching test-policy-plr with decision of cluster managed1")
			plr := utils.GetWithTimeout(clientHubDynamic, gvrPlacementRule, case1PolicyName+"-plr", testNamespace, true, defaultTimeoutSeconds)
			plr.Object["status"] = utils.GeneratePlrStatus("managed1")
			plr, err := clientHubDynamic.Resource(gvrPlacementRule).Namespace(testNamespace).UpdateStatus(plr, metav1.UpdateOptions{})
			Expect(err).To(BeNil())
			plc := utils.GetWithTimeout(clientHubDynamic, gvrPolicy, testNamespace+"."+case1PolicyName, "managed1", true, defaultTimeoutSeconds)
			Expect(plc).ToNot(BeNil())
			opt := metav1.ListOptions{LabelSelector: common.RootPolicyLabel + "=" + testNamespace + "." + case1PolicyName}
			utils.ListWithTimeout(clientHubDynamic, gvrPolicy, opt, 1, true, defaultTimeoutSeconds)
		})
		It("should update replicated policy in ns managed1", func() {
			By("Patching test-policy with spec.remediationAction = enforce")
			rootPlc := utils.GetWithTimeout(clientHubDynamic, gvrPolicy, case1PolicyName, testNamespace, true, defaultTimeoutSeconds)
			Expect(rootPlc).NotTo(BeNil())
			Expect(rootPlc.Object["spec"].(map[string]interface{})["remediationAction"]).To(Equal("inform"))
			rootPlc.Object["spec"].(map[string]interface{})["remediationAction"] = "enforce"
			rootPlc, err := clientHubDynamic.Resource(gvrPolicy).Namespace(testNamespace).Update(rootPlc, metav1.UpdateOptions{})
			Expect(err).To(BeNil())
			Eventually(func() interface{} {
				replicatedPlc := utils.GetWithTimeout(clientHubDynamic, gvrPolicy, testNamespace+"."+case1PolicyName, "managed1", true, defaultTimeoutSeconds)
				return replicatedPlc.Object["spec"]
			}, defaultTimeoutSeconds, 1).Should(utils.SemanticEqual(rootPlc.Object["spec"]))
		})
		It("should remove replicated policy in ns managed1", func() {
			By("Patching test-policy with spec.disabled = true")
			rootPlc := utils.GetWithTimeout(clientHubDynamic, gvrPolicy, case1PolicyName, testNamespace, true, defaultTimeoutSeconds)
			Expect(rootPlc).NotTo(BeNil())
			Expect(rootPlc.Object["spec"].(map[string]interface{})["disabled"]).To(Equal(false))
			rootPlc.Object["spec"].(map[string]interface{})["disabled"] = true
			rootPlc, err := clientHubDynamic.Resource(gvrPolicy).Namespace(testNamespace).Update(rootPlc, metav1.UpdateOptions{})
			Expect(err).To(BeNil())
			Expect(rootPlc.Object["spec"].(map[string]interface{})["disabled"]).To(Equal(true))
			opt := metav1.ListOptions{LabelSelector: common.RootPolicyLabel + "=" + testNamespace + "." + case1PolicyName}
			utils.ListWithTimeout(clientHubDynamic, gvrPolicy, opt, 0, true, defaultTimeoutSeconds)
		})
		It("should be created in user ns", func() {
			By("Creating " + case1PolicyYaml)
			utils.Kubectl("apply",
				"-f", case1PolicyYaml,
				"-n", testNamespace)
			plc := utils.GetWithTimeout(clientHubDynamic, gvrPolicy, case1PolicyName, testNamespace, true, defaultTimeoutSeconds)
			Expect(plc).NotTo(BeNil())
		})
		It("should propagate to cluster ns managed1", func() {
			By("Patching test-policy-plr with decision of cluster managed1")
			plr := utils.GetWithTimeout(clientHubDynamic, gvrPlacementRule, case1PolicyName+"-plr", testNamespace, true, defaultTimeoutSeconds)
			plr.Object["status"] = utils.GeneratePlrStatus("managed1")
			plr, err := clientHubDynamic.Resource(gvrPlacementRule).Namespace(testNamespace).UpdateStatus(plr, metav1.UpdateOptions{})
			Expect(err).To(BeNil())
			plc := utils.GetWithTimeout(clientHubDynamic, gvrPolicy, testNamespace+"."+case1PolicyName, "managed1", true, defaultTimeoutSeconds)
			Expect(plc).ToNot(BeNil())
			opt := metav1.ListOptions{LabelSelector: common.RootPolicyLabel + "=" + testNamespace + "." + case1PolicyName}
			utils.ListWithTimeout(clientHubDynamic, gvrPolicy, opt, 1, true, defaultTimeoutSeconds)
		})
		It("should update test-policy to a different policy template", func() {
			By("Creating ../resources/case1_propagation/case1-test-policy2.yaml")
			utils.Kubectl("apply",
				"-f", "../resources/case1_propagation/case1-test-policy2.yaml",
				"-n", testNamespace)
			rootPlc := utils.GetWithTimeout(clientHubDynamic, gvrPolicy, case1PolicyName, testNamespace, true, defaultTimeoutSeconds)
			Expect(rootPlc).NotTo(BeNil())
			yamlPlc := utils.ParseYaml("../resources/case1_propagation/case1-test-policy2.yaml")
			Eventually(func() interface{} {
				replicatedPlc := utils.GetWithTimeout(clientHubDynamic, gvrPolicy, testNamespace+"."+case1PolicyName, "managed1", true, defaultTimeoutSeconds)
				return replicatedPlc.Object["spec"]
			}, defaultTimeoutSeconds, 1).Should(utils.SemanticEqual(yamlPlc.Object["spec"]))
		})
		It("should clean up", func() {
			utils.Kubectl("delete",
				"-f", "../resources/case1_propagation/case1-test-policy2.yaml",
				"-n", testNamespace)
			opt := metav1.ListOptions{}
			utils.ListWithTimeout(clientHubDynamic, gvrPolicy, opt, 0, false, 10)
		})
	})
})