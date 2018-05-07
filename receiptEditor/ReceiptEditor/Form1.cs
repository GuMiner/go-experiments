using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Data;
using System.Drawing;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Windows.Forms;

namespace ReceiptEditor
{
    public partial class ReceiptEditor : Form
    {
        public ReceiptEditor()
        {
            InitializeComponent();
            this.imageBox.Image = Image.FromFile(@"C:\Users\Gustave\Desktop\Data Archive\receipts\2015-01-21\001.jpg");
        }

        // Unused?
        private void Form1_Load(object sender, EventArgs e)
        {

        }

        private void nextButton_Click(object sender, EventArgs e)
        {

        }

        // Unused
        private void label2_Click(object sender, EventArgs e)
        {
            
        }

        private void textBox1_TextChanged(object sender, EventArgs e)
        {

        }

        private void categoryBox_SelectedIndexChanged(object sender, EventArgs e)
        {

        }

        private void itemDate_ValueChanged(object sender, EventArgs e)
        {

        }

        private void imageBox_Click(object sender, EventArgs e)
        {

        }

        private void imageBox_Paint(object sender, PaintEventArgs e)
        {

        }

        private void imageBox_MouseMove(object sender, MouseEventArgs e)
        {

        }

        private void imageBox_MouseDown(object sender, MouseEventArgs e)
        {

        }

        private void imageBox_MouseUp(object sender, MouseEventArgs e)
        {

        }

        private void button1_Click(object sender, EventArgs e)
        {

        }

        private void addCategoryField_TextChanged(object sender, EventArgs e)
        {

        }
    }
}
